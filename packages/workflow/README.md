# @labhaus/workflow

通用工作流引擎，提供类型安全的状态机和工作流执行能力。

## 核心概念

### 工作流定义 (WorkflowDefinition)

工作流由**节点**和**边**组成：

- **节点 (Node)**: 执行单元，有输入输出
- **边 (Edge)**: 连接节点，定义数据流

### 工作流执行 (WorkflowExecution)

工作流的执行实例，包含：

- 执行状态（draft, pending, running, paused, completed, failed, cancelled）
- 上下文数据（节点间共享）
- 审计事件

## 使用示例

### 1. 定义节点处理器

```typescript
import { NodeHandler, NodeType, WorkflowNode } from '@labhaus/workflow';

const inputHandler: NodeHandler = {
  type: NodeType.INPUT,
  async execute(node: WorkflowNode, context: Record<string, unknown>) {
    // 从外部读取数据
    return { input_data: 'Hello' };
  },
};

const processHandler: NodeHandler = {
  type: NodeType.PROCESS,
  async execute(node: WorkflowNode, context: Record<string, unknown>) {
    const input = context.input_data as string;
    return { output_data: input.toUpperCase() };
  },
};
```

### 2. 创建工作流

```typescript
import { WorkflowDefinition, WorkflowExecution, WorkflowExecutor } from '@labhaus/workflow';

const workflow: WorkflowDefinition = {
  id: crypto.randomUUID(),
  name: 'Simple Text Transform',
  nodes: [
    { id: 'input', type: NodeType.INPUT, name: 'Input', inputs: [], outputs: [] },
    { id: 'process', type: NodeType.PROCESS, name: 'Process', inputs: [], outputs: [] },
  ],
  edges: [{ id: 'edge1', source: 'input', target: 'process' }],
  version: 1,
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
};
```

### 3. 执行工作流

```typescript
const executor = new WorkflowExecutor();
executor.registerHandler(inputHandler);
executor.registerHandler(processHandler);

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

console.log(result.status); // 'completed'
console.log(result.context); // { input_data: 'Hello', output_data: 'HELLO' }
```

## 状态机

### 工作流状态转换

```
DRAFT → PENDING → RUNNING → COMPLETED
                     ↓
                  PAUSED → RUNNING
                     ↓
                  FAILED → PENDING (retry)
                     ↓
                 CANCELLED
```

### 合法转换

- `DRAFT` → `PENDING`
- `PENDING` → `RUNNING` | `CANCELLED`
- `RUNNING` → `PAUSED` | `COMPLETED` | `FAILED` | `CANCELLED`
- `PAUSED` → `RUNNING` | `CANCELLED`
- `FAILED` → `PENDING` (允许重试)
- `COMPLETED`, `CANCELLED` 为终态

### 状态检查

```typescript
import { canTransition, transitionWorkflowStatus } from '@labhaus/workflow';

// 检查是否可以转换
if (canTransition(WorkflowStatus.DRAFT, WorkflowStatus.PENDING)) {
  const result = transitionWorkflowStatus(
    execution.id,
    WorkflowStatus.DRAFT,
    WorkflowStatus.PENDING,
    'user-123',
    'Starting workflow'
  );
  
  console.log(result.event); // 审计事件
}
```

## 节点类型

- `INPUT`: 输入节点（数据源）
- `PROCESS`: 处理节点（业务逻辑）
- `OUTPUT`: 输出节点（数据目的地）
- `CONDITION`: 条件节点（分支判断）
- `APPROVAL`: 审批节点（人工介入）

## 工作流验证

```typescript
const validation = executor.validateWorkflow(workflow);

if (!validation.valid) {
  console.error('Validation errors:', validation.errors);
}
```

验证规则：

- 至少有一个节点
- 所有节点类型都有注册的处理器
- 边引用的节点必须存在
- 不允许循环依赖

## API 参考

### WorkflowExecutor

**方法：**

- `registerHandler(handler: NodeHandler)` - 注册节点处理器
- `validateWorkflow(workflow: WorkflowDefinition)` - 验证工作流定义
- `executeWorkflow(workflow, execution, actor?)` - 执行工作流

### 状态机函数

- `canTransition(current, target)` - 检查是否可转换
- `transitionWorkflowStatus(executionId, current, target, actor, note?)` - 执行状态转换
- `getNextStatuses(current)` - 获取允许的下一状态
- `isTerminalStatus(status)` - 检查是否为终态

## 测试

```bash
pnpm test
```

## 设计来源

基于 [ai-video-factory](https://github.com/sine-io/ai-video-factory) 的任务状态机重构而来。
