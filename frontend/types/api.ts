// API Response types
export interface ApiResponse<T = unknown> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

// User types
export interface User {
  id: string;
  email: string;
  username: string;
  phone?: string;
  avatar_url?: string;
  role: 'user' | 'vip' | 'admin';
  balance: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
}

export interface AuthResponse {
  user: User;
  token: string;
  refresh_token: string;
  expires_in: number;
}

// Agent types
export interface Agent {
  id: string;
  agent_id: string;
  name: string;
  name_en?: string;
  description: string;
  category: string;
  tags: string[];
  tier: 1 | 2 | 3 | 4;
  runtime_type: 'openclaw' | 'claude' | 'openai';
  price_per_request: number;
  price_per_token?: number;
  avg_duration_seconds: number;
  success_rate: number;
  rating: number;
  total_tasks: number;
  total_success: number;
  total_failed: number;
  is_active: boolean;
  input_example?: string;
  output_example?: string;
  created_at: string;
  updated_at: string;
}

export interface AgentListResponse {
  total: number;
  page: number;
  limit: number;
  agents: Agent[];
}

// Task types
export type TaskStatus = 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled';

export interface Task {
  id: string;
  user_id: string;
  agent_id: string;
  agent?: Agent;
  prompt: string;
  result?: string;
  status: TaskStatus;
  priority: number;
  cost: number;
  tokens_used: number;
  duration_seconds: number;
  error_message?: string;
  retry_count: number;
  max_retries: number;
  created_at: string;
  started_at?: string;
  completed_at?: string;
}

export interface CreateTaskRequest {
  agent_id: string;
  prompt: string;
  priority?: number;
  max_retries?: number;
}

export interface CreateTaskResponse {
  task_id: string;
  status: TaskStatus;
  estimated_cost: number;
  estimated_duration: number;
}

export interface TaskListResponse {
  total: number;
  page: number;
  limit: number;
  tasks: Task[];
}

// Plan and Subscription types
export type PlanType = 'basic' | 'pro' | 'enterprise';
export type BillingCycle = 'monthly' | 'yearly';
export type TransactionType = 'recharge' | 'consumption' | 'refund' | 'subscription';

export interface Plan {
  id: string;
  name: string;
  type: PlanType;
  price_monthly: number;
  price_yearly: number;
  features: {
    task_limit: number | null;  // null = unlimited
    agent_tier_limit: number;
    api_access: boolean;
    priority_support: boolean;
    custom_agents: boolean;
  };
  is_active: boolean;
}

export interface Subscription {
  id: string;
  user_id: string;
  plan_id: string;
  plan_type: PlanType;
  status: 'active' | 'cancelled' | 'expired';
  start_date: string;
  end_date: string;
  price: number;
  billing_cycle: BillingCycle;
  auto_renew: boolean;
  created_at: string;
}

export interface Transaction {
  id: string;
  user_id: string;
  type: TransactionType;
  amount: number;
  balance_before: number;
  balance_after: number;
  description: string;
  reference_id?: string;
  created_at: string;
}

export interface UserBalance {
  balance: number;
  subscription?: Subscription;
}

// Category types
export interface Category {
  id: string;
  name: string;
  slug: string;
  icon: string;
  description: string;
  sort_order: number;
  agent_count: number;
}

// Pagination
export interface Pagination {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

// API Error
export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, string[]>;
}
