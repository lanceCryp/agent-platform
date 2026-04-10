import axios, { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from 'axios';
import type { ApiResponse, ApiError } from '@/types/api';

// API Base URL - in production, this would be from environment
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// Create axios instance
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Token management
let accessToken: string | null = null;
let refreshToken: string | null = null;

export const setTokens = (access: string, refresh: string) => {
  accessToken = access;
  refreshToken = refresh;
  if (typeof window !== 'undefined') {
    localStorage.setItem('access_token', access);
    localStorage.setItem('refresh_token', refresh);
  }
};

export const clearTokens = () => {
  accessToken = null;
  refreshToken = null;
  if (typeof window !== 'undefined') {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
  }
};

export const getAccessToken = () => {
  if (accessToken) return accessToken;
  if (typeof window !== 'undefined') {
    accessToken = localStorage.getItem('access_token');
  }
  return accessToken;
};

// Request interceptor - add auth token
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getAccessToken();
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor - handle errors
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ApiError>) => {
    const originalRequest = error.config;

    // Handle 401 - try to refresh token
    if (error.response?.status === 401 && refreshToken && originalRequest) {
      try {
        const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
          refresh_token: refreshToken,
        });

        const { token: newAccessToken, refresh_token: newRefreshToken } = response.data.data;
        setTokens(newAccessToken, newRefreshToken);

        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
        }
        return apiClient(originalRequest);
      } catch (refreshError) {
        clearTokens();
        // Redirect to login
        if (typeof window !== 'undefined') {
          window.location.href = '/login';
        }
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

// API methods
export const api = {
  // Auth
  auth: {
    login: (email: string, password: string) =>
      apiClient.post<ApiResponse<{ token: string; refresh_token: string; user: unknown }>>('/auth/login', {
        email,
        password,
      }),
    register: (email: string, username: string, password: string) =>
      apiClient.post<ApiResponse<{ token: string; refresh_token: string; user: unknown }>>('/auth/register', {
        email,
        username,
        password,
      }),
    logout: () => apiClient.post('/auth/logout'),
    refresh: (refreshToken: string) =>
      apiClient.post('/auth/refresh', { refresh_token: refreshToken }),
    me: () => apiClient.get('/users/me'),
  },

  // Users
  users: {
    getProfile: () => apiClient.get('/users/me'),
    updateProfile: (data: { username?: string; phone?: string; avatar_url?: string }) =>
      apiClient.patch('/users/me', data),
    getBalance: () => apiClient.get<ApiResponse<{ balance: number; subscription?: unknown }>>('/users/me/balance'),
  },

  // Agents
  agents: {
    list: (params?: { category?: string; tier?: number; page?: number; limit?: number; search?: string }) =>
      apiClient.get<ApiResponse<{ total: number; agents: unknown[] }>>('/agents', { params }),
    get: (agentId: string) => apiClient.get<ApiResponse<unknown>>(`/agents/${agentId}`),
    getCategories: () => apiClient.get<ApiResponse<unknown[]>>('/agents/categories'),
  },

  // Tasks
  tasks: {
    create: (data: { agent_id: string; prompt: string; priority?: number; max_retries?: number }) =>
      apiClient.post<ApiResponse<{ task_id: string; status: string; estimated_cost: number }>>('/tasks', data),
    get: (taskId: string) => apiClient.get<ApiResponse<unknown>>(`/tasks/${taskId}`),
    list: (params?: { status?: string; page?: number; limit?: number }) =>
      apiClient.get<ApiResponse<{ total: number; tasks: unknown[] }>>('/tasks', { params }),
    cancel: (taskId: string) => apiClient.post(`/tasks/${taskId}/cancel`),
    retry: (taskId: string) => apiClient.post(`/tasks/${taskId}/retry`),
  },

  // Subscriptions
  subscriptions: {
    list: () => apiClient.get<ApiResponse<unknown[]>>('/subscriptions'),
    create: (planId: string, billingCycle: 'monthly' | 'yearly') =>
      apiClient.post<ApiResponse<{ subscription_id: string; payment_url: string }>>('/subscriptions', {
        plan_id: planId,
        billing_cycle: billingCycle,
      }),
    cancel: (subscriptionId: string) => apiClient.post(`/subscriptions/${subscriptionId}/cancel`),
  },

  // Plans
  plans: {
    list: () => apiClient.get<ApiResponse<unknown[]>>('/plans'),
    get: (planId: string) => apiClient.get<ApiResponse<unknown>>(`/plans/${planId}`),
  },

  // Transactions
  transactions: {
    list: (params?: { page?: number; limit?: number; type?: string }) =>
      apiClient.get<ApiResponse<unknown[]>>('/transactions', { params }),
  },

  // Billing
  billing: {
    recharge: (amount: number, method: string) =>
      apiClient.post<ApiResponse<{ payment_url: string }>>('/billing/recharge', {
        amount,
        method,
      }),
    getHistory: (params?: { page?: number; limit?: number }) =>
      apiClient.get<ApiResponse<unknown[]>>('/billing/history', { params }),
  },
};

export default api;
