import { z } from 'zod';

/**
 * Workflow execution status
 */
export enum WorkflowStatus {
  DRAFT = 'draft',
  PENDING = 'pending',
  RUNNING = 'running',
  PAUSED = 'paused',
  COMPLETED = 'completed',
  FAILED = 'failed',
  CANCELLED = 'cancelled',
}

/**
 * Node execution status
 */
export enum NodeStatus {
  PENDING = 'pending',
  RUNNING = 'running',
  COMPLETED = 'completed',
  FAILED = 'failed',
  SKIPPED = 'skipped',
}

/**
 * Node types
 */
export enum NodeType {
  INPUT = 'input',
  PROCESS = 'process',
  OUTPUT = 'output',
  CONDITION = 'condition',
  APPROVAL = 'approval',
}

/**
 * Workflow node definition
 */
export const workflowNodeSchema = z.object({
  id: z.string(),
  type: z.nativeEnum(NodeType),
  name: z.string(),
  config: z.record(z.unknown()).optional(),
  inputs: z.array(z.string()).default([]),
  outputs: z.array(z.string()).default([]),
});

/**
 * Workflow edge (connection between nodes)
 */
export const workflowEdgeSchema = z.object({
  id: z.string(),
  source: z.string(),
  target: z.string(),
  condition: z.string().optional(),
});

/**
 * Workflow definition
 */
export const workflowDefinitionSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  description: z.string().optional(),
  nodes: z.array(workflowNodeSchema),
  edges: z.array(workflowEdgeSchema),
  version: z.number().int().positive().default(1),
  created_at: z.string().datetime(),
  updated_at: z.string().datetime(),
});

/**
 * Workflow execution instance
 */
export const workflowExecutionSchema = z.object({
  id: z.string().uuid(),
  workflow_id: z.string().uuid(),
  status: z.nativeEnum(WorkflowStatus),
  current_node: z.string().nullable(),
  context: z.record(z.unknown()).default({}),
  error: z.string().nullable().optional(),
  started_at: z.string().datetime().nullable(),
  completed_at: z.string().datetime().nullable(),
  created_at: z.string().datetime(),
  updated_at: z.string().datetime(),
});

/**
 * Workflow transition event
 */
export const workflowTransitionEventSchema = z.object({
  execution_id: z.string().uuid(),
  event_type: z.string(),
  actor: z.string(),
  from_status: z.nativeEnum(WorkflowStatus),
  to_status: z.nativeEnum(WorkflowStatus),
  note: z.string().optional(),
  created_at: z.string().datetime(),
});

export type WorkflowNode = z.infer<typeof workflowNodeSchema>;
export type WorkflowEdge = z.infer<typeof workflowEdgeSchema>;
export type WorkflowDefinition = z.infer<typeof workflowDefinitionSchema>;
export type WorkflowExecution = z.infer<typeof workflowExecutionSchema>;
export type WorkflowTransitionEvent = z.infer<typeof workflowTransitionEventSchema>;

/**
 * Node handler interface
 */
export interface NodeHandler {
  type: NodeType;
  execute(node: WorkflowNode, context: Record<string, unknown>): Promise<Record<string, unknown>>;
  validate?(node: WorkflowNode): boolean;
}
