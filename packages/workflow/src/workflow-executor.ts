import {
  WorkflowDefinition,
  WorkflowExecution,
  WorkflowNode,
  WorkflowEdge,
  NodeHandler,
  WorkflowStatus,
  NodeStatus,
} from './types.js';
import { transitionWorkflowStatus } from './workflow-state.js';

/**
 * Workflow execution error
 */
export class WorkflowExecutionError extends Error {
  constructor(
    message: string,
    public executionId: string,
    public nodeId?: string
  ) {
    super(message);
    this.name = 'WorkflowExecutionError';
  }
}

/**
 * Workflow executor
 */
export class WorkflowExecutor {
  private handlers = new Map<string, NodeHandler>();

  /**
   * Register a node handler
   */
  registerHandler(handler: NodeHandler): void {
    this.handlers.set(handler.type, handler);
  }

  /**
   * Get handler for node type
   */
  private getHandler(node: WorkflowNode): NodeHandler {
    const handler = this.handlers.get(node.type);
    if (!handler) {
      throw new WorkflowExecutionError(
        `No handler registered for node type: ${node.type}`,
        '',
        node.id
      );
    }
    return handler;
  }

  /**
   * Validate workflow definition
   */
  validateWorkflow(workflow: WorkflowDefinition): { valid: boolean; errors: string[] } {
    const errors: string[] = [];

    // Check for nodes
    if (workflow.nodes.length === 0) {
      errors.push('Workflow must have at least one node');
    }

    // Check for cycles (simple DFS)
    const visited = new Set<string>();
    const recursionStack = new Set<string>();

    const hasCycle = (nodeId: string): boolean => {
      visited.add(nodeId);
      recursionStack.add(nodeId);

      const outgoingEdges = workflow.edges.filter((e) => e.source === nodeId);
      for (const edge of outgoingEdges) {
        if (!visited.has(edge.target)) {
          if (hasCycle(edge.target)) return true;
        } else if (recursionStack.has(edge.target)) {
          return true;
        }
      }

      recursionStack.delete(nodeId);
      return false;
    };

    for (const node of workflow.nodes) {
      if (!visited.has(node.id) && hasCycle(node.id)) {
        errors.push(`Workflow contains cycle involving node: ${node.id}`);
        break;
      }
    }

    // Validate each node has a registered handler
    for (const node of workflow.nodes) {
      if (!this.handlers.has(node.type)) {
        errors.push(`No handler registered for node type: ${node.type} (node: ${node.id})`);
      }
    }

    // Validate edges reference existing nodes
    const nodeIds = new Set(workflow.nodes.map((n) => n.id));
    for (const edge of workflow.edges) {
      if (!nodeIds.has(edge.source)) {
        errors.push(`Edge ${edge.id} references non-existent source node: ${edge.source}`);
      }
      if (!nodeIds.has(edge.target)) {
        errors.push(`Edge ${edge.id} references non-existent target node: ${edge.target}`);
      }
    }

    return {
      valid: errors.length === 0,
      errors,
    };
  }

  /**
   * Find start nodes (nodes with no incoming edges)
   */
  private findStartNodes(workflow: WorkflowDefinition): WorkflowNode[] {
    const nodesWithIncoming = new Set(workflow.edges.map((e) => e.target));
    return workflow.nodes.filter((n) => !nodesWithIncoming.has(n.id));
  }

  /**
   * Find next nodes to execute
   */
  private findNextNodes(
    workflow: WorkflowDefinition,
    currentNode: WorkflowNode,
    context: Record<string, unknown>
  ): WorkflowNode[] {
    const outgoingEdges = workflow.edges.filter((e) => e.source === currentNode.id);
    const nextNodeIds: string[] = [];

    for (const edge of outgoingEdges) {
      // If edge has a condition, evaluate it
      if (edge.condition) {
        // Simple condition evaluation: check if context[condition] is truthy
        if (context[edge.condition]) {
          nextNodeIds.push(edge.target);
        }
      } else {
        nextNodeIds.push(edge.target);
      }
    }

    return workflow.nodes.filter((n) => nextNodeIds.includes(n.id));
  }

  /**
   * Execute a single node
   */
  async executeNode(
    node: WorkflowNode,
    context: Record<string, unknown>
  ): Promise<Record<string, unknown>> {
    const handler = this.getHandler(node);

    // Validate node if handler provides validation
    if (handler.validate && !handler.validate(node)) {
      throw new WorkflowExecutionError(
        `Node validation failed: ${node.id}`,
        '',
        node.id
      );
    }

    // Execute node
    return await handler.execute(node, context);
  }

  /**
   * Execute workflow (simplified synchronous execution)
   */
  async executeWorkflow(
    workflow: WorkflowDefinition,
    execution: WorkflowExecution,
    actor: string = 'system'
  ): Promise<WorkflowExecution> {
    // Validate workflow
    const validation = this.validateWorkflow(workflow);
    if (!validation.valid) {
      throw new WorkflowExecutionError(
        `Workflow validation failed: ${validation.errors.join(', ')}`,
        execution.id
      );
    }

    // Transition draft -> pending first if needed
    if (execution.status === WorkflowStatus.DRAFT) {
      transitionWorkflowStatus(execution.id, execution.status, WorkflowStatus.PENDING, actor);
      execution.status = WorkflowStatus.PENDING;
    }

    // Transition to running
    transitionWorkflowStatus(execution.id, execution.status, WorkflowStatus.RUNNING, actor);
    execution.status = WorkflowStatus.RUNNING;
    execution.started_at = new Date().toISOString();

    const context = { ...execution.context };

    try {
      // Find start nodes
      const startNodes = this.findStartNodes(workflow);
      if (startNodes.length === 0) {
        throw new WorkflowExecutionError('No start nodes found in workflow', execution.id);
      }

      // Simple sequential execution (BFS)
      const queue = [...startNodes];
      const executed = new Set<string>();

      while (queue.length > 0) {
        const node = queue.shift()!;
        if (executed.has(node.id)) continue;

        // Execute node
        const result = await this.executeNode(node, context);
        Object.assign(context, result);
        executed.add(node.id);

        // Find next nodes
        const nextNodes = this.findNextNodes(workflow, node, context);
        queue.push(...nextNodes);
      }

      // Success
      execution.status = WorkflowStatus.COMPLETED;
      execution.completed_at = new Date().toISOString();
      execution.context = context;
    } catch (error) {
      // Failure
      execution.status = WorkflowStatus.FAILED;
      execution.error = error instanceof Error ? error.message : String(error);
      execution.completed_at = new Date().toISOString();
      throw error;
    }

    return execution;
  }
}
