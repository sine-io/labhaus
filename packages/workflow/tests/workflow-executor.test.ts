import { describe, it, expect, beforeEach } from 'vitest';
import { WorkflowExecutor, WorkflowExecutionError } from '../src/workflow-executor';
import {
  WorkflowDefinition,
  WorkflowExecution,
  WorkflowStatus,
  NodeType,
  NodeHandler,
} from '../src/types';

describe('WorkflowExecutor', () => {
  let executor: WorkflowExecutor;

  beforeEach(() => {
    executor = new WorkflowExecutor();
  });

  describe('validateWorkflow', () => {
    it('should validate a simple valid workflow', () => {
      const mockHandler: NodeHandler = {
        type: NodeType.INPUT,
        execute: async () => ({}),
      };
      executor.registerHandler(mockHandler);

      const workflow: WorkflowDefinition = {
        id: crypto.randomUUID(),
        name: 'Test Workflow',
        nodes: [
          { id: 'node1', type: NodeType.INPUT, name: 'Input', inputs: [], outputs: [] },
        ],
        edges: [],
        version: 1,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      const result = executor.validateWorkflow(workflow);
      expect(result.valid).toBe(true);
      expect(result.errors).toHaveLength(0);
    });

    it('should reject workflow with no nodes', () => {
      const workflow: WorkflowDefinition = {
        id: crypto.randomUUID(),
        name: 'Empty Workflow',
        nodes: [],
        edges: [],
        version: 1,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      const result = executor.validateWorkflow(workflow);
      expect(result.valid).toBe(false);
      expect(result.errors).toContain('Workflow must have at least one node');
    });

    it('should reject workflow with unregistered handler', () => {
      const workflow: WorkflowDefinition = {
        id: crypto.randomUUID(),
        name: 'Test Workflow',
        nodes: [
          { id: 'node1', type: NodeType.INPUT, name: 'Input', inputs: [], outputs: [] },
        ],
        edges: [],
        version: 1,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      const result = executor.validateWorkflow(workflow);
      expect(result.valid).toBe(false);
      expect(result.errors.some((e) => e.includes('No handler registered'))).toBe(true);
    });

    it('should detect cycles', () => {
      const mockHandler: NodeHandler = {
        type: NodeType.PROCESS,
        execute: async () => ({}),
      };
      executor.registerHandler(mockHandler);

      const workflow: WorkflowDefinition = {
        id: crypto.randomUUID(),
        name: 'Cyclic Workflow',
        nodes: [
          { id: 'node1', type: NodeType.PROCESS, name: 'Node 1', inputs: [], outputs: [] },
          { id: 'node2', type: NodeType.PROCESS, name: 'Node 2', inputs: [], outputs: [] },
        ],
        edges: [
          { id: 'edge1', source: 'node1', target: 'node2' },
          { id: 'edge2', source: 'node2', target: 'node1' }, // Cycle
        ],
        version: 1,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      const result = executor.validateWorkflow(workflow);
      expect(result.valid).toBe(false);
      expect(result.errors.some((e) => e.includes('cycle'))).toBe(true);
    });
  });

  describe('executeWorkflow', () => {
    it('should execute a simple workflow', async () => {
      const inputHandler: NodeHandler = {
        type: NodeType.INPUT,
        execute: async (node, context) => ({ input_value: 'test' }),
      };

      const processHandler: NodeHandler = {
        type: NodeType.PROCESS,
        execute: async (node, context) => ({
          output_value: `processed-${context.input_value}`,
        }),
      };

      executor.registerHandler(inputHandler);
      executor.registerHandler(processHandler);

      const workflow: WorkflowDefinition = {
        id: crypto.randomUUID(),
        name: 'Simple Workflow',
        nodes: [
          { id: 'input', type: NodeType.INPUT, name: 'Input', inputs: [], outputs: [] },
          { id: 'process', type: NodeType.PROCESS, name: 'Process', inputs: [], outputs: [] },
        ],
        edges: [{ id: 'edge1', source: 'input', target: 'process' }],
        version: 1,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      const execution: WorkflowExecution = {
        id: crypto.randomUUID(),
        workflow_id: workflow.id,
        status: WorkflowStatus.DRAFT,
        current_node: null,
        context: {},
        started_at: null,
        completed_at: null,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      const result = await executor.executeWorkflow(workflow, execution);

      expect(result.status).toBe(WorkflowStatus.COMPLETED);
      expect(result.context.input_value).toBe('test');
      expect(result.context.output_value).toBe('processed-test');
      expect(result.completed_at).not.toBeNull();
    });

    it('should handle execution errors', async () => {
      const failingHandler: NodeHandler = {
        type: NodeType.PROCESS,
        execute: async () => {
          throw new Error('Execution failed');
        },
      };

      executor.registerHandler(failingHandler);

      const workflow: WorkflowDefinition = {
        id: crypto.randomUUID(),
        name: 'Failing Workflow',
        nodes: [
          { id: 'fail', type: NodeType.PROCESS, name: 'Fail', inputs: [], outputs: [] },
        ],
        edges: [],
        version: 1,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      const execution: WorkflowExecution = {
        id: crypto.randomUUID(),
        workflow_id: workflow.id,
        status: WorkflowStatus.DRAFT,
        current_node: null,
        context: {},
        started_at: null,
        completed_at: null,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      await expect(executor.executeWorkflow(workflow, execution)).rejects.toThrow(
        'Execution failed'
      );
      expect(execution.status).toBe(WorkflowStatus.FAILED);
      expect(execution.error).toContain('Execution failed');
    });
  });
});
