import { describe, it, expect } from 'vitest';
import {
  canTransition,
  transitionWorkflowStatus,
  getNextStatuses,
  isTerminalStatus,
  InvalidStatusTransition,
} from '../src/workflow-state';
import { WorkflowStatus } from '../src/types';

describe('WorkflowState', () => {
  describe('canTransition', () => {
    it('should allow draft -> pending', () => {
      expect(canTransition(WorkflowStatus.DRAFT, WorkflowStatus.PENDING)).toBe(true);
    });

    it('should allow pending -> running', () => {
      expect(canTransition(WorkflowStatus.PENDING, WorkflowStatus.RUNNING)).toBe(true);
    });

    it('should allow running -> completed', () => {
      expect(canTransition(WorkflowStatus.RUNNING, WorkflowStatus.COMPLETED)).toBe(true);
    });

    it('should allow running -> failed', () => {
      expect(canTransition(WorkflowStatus.RUNNING, WorkflowStatus.FAILED)).toBe(true);
    });

    it('should allow failed -> pending (retry)', () => {
      expect(canTransition(WorkflowStatus.FAILED, WorkflowStatus.PENDING)).toBe(true);
    });

    it('should not allow completed -> running', () => {
      expect(canTransition(WorkflowStatus.COMPLETED, WorkflowStatus.RUNNING)).toBe(false);
    });

    it('should not allow draft -> running', () => {
      expect(canTransition(WorkflowStatus.DRAFT, WorkflowStatus.RUNNING)).toBe(false);
    });

    it('should allow running -> paused', () => {
      expect(canTransition(WorkflowStatus.RUNNING, WorkflowStatus.PAUSED)).toBe(true);
    });

    it('should allow paused -> running', () => {
      expect(canTransition(WorkflowStatus.PAUSED, WorkflowStatus.RUNNING)).toBe(true);
    });
  });

  describe('transitionWorkflowStatus', () => {
    it('should successfully transition with valid states', () => {
      const result = transitionWorkflowStatus(
        'exec-123',
        WorkflowStatus.DRAFT,
        WorkflowStatus.PENDING,
        'user-1',
        'Starting workflow'
      );

      expect(result.status).toBe(WorkflowStatus.PENDING);
      expect(result.event.execution_id).toBe('exec-123');
      expect(result.event.from_status).toBe(WorkflowStatus.DRAFT);
      expect(result.event.to_status).toBe(WorkflowStatus.PENDING);
      expect(result.event.actor).toBe('user-1');
      expect(result.event.note).toBe('Starting workflow');
    });

    it('should throw on invalid transition', () => {
      expect(() =>
        transitionWorkflowStatus(
          'exec-123',
          WorkflowStatus.DRAFT,
          WorkflowStatus.COMPLETED,
          'user-1'
        )
      ).toThrow(InvalidStatusTransition);
    });
  });

  describe('getNextStatuses', () => {
    it('should return allowed next statuses for draft', () => {
      const next = getNextStatuses(WorkflowStatus.DRAFT);
      expect(next).toContain(WorkflowStatus.PENDING);
      expect(next).toHaveLength(1);
    });

    it('should return empty array for completed', () => {
      const next = getNextStatuses(WorkflowStatus.COMPLETED);
      expect(next).toHaveLength(0);
    });
  });

  describe('isTerminalStatus', () => {
    it('should return true for completed', () => {
      expect(isTerminalStatus(WorkflowStatus.COMPLETED)).toBe(true);
    });

    it('should return true for cancelled', () => {
      expect(isTerminalStatus(WorkflowStatus.CANCELLED)).toBe(true);
    });

    it('should return false for running', () => {
      expect(isTerminalStatus(WorkflowStatus.RUNNING)).toBe(false);
    });
  });
});
