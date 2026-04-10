"use client";

import * as React from "react";
import { clsx } from "clsx";

// Toast types
type ToastType = "success" | "error" | "warning" | "info";

interface Toast {
  id: string;
  type: ToastType;
  title: string;
  description?: string;
  duration?: number;
}

// Toast context
interface ToastContextType {
  toasts: Toast[];
  addToast: (toast: Omit<Toast, "id">) => void;
  removeToast: (id: string) => void;
}

const ToastContext = React.createContext<ToastContextType | undefined>(undefined);

// Provider component
export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = React.useState<Toast[]>([]);

  const addToast = React.useCallback((toast: Omit<Toast, "id">) => {
    const id = Math.random().toString(36).substring(2, 9);
    const newToast = { ...toast, id };
    
    setToasts((prev) => [...prev, newToast]);

    // Auto remove after duration
    const duration = toast.duration || 5000;
    setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id));
    }, duration);
  }, []);

  const removeToast = React.useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  return (
    <ToastContext.Provider value={{ toasts, addToast, removeToast }}>
      {children}
      <ToastContainer />
    </ToastContext.Provider>
  );
}

// useToast hook
export function useToast() {
  const context = React.useContext(ToastContext);
  if (!context) {
    throw new Error("useToast must be used within a ToastProvider");
  }

  const toast = {
    success: (title: string, description?: string) => {
      context.addToast({ type: "success", title, description });
    },
    error: (title: string, description?: string) => {
      context.addToast({ type: "error", title, description, duration: 8000 });
    },
    warning: (title: string, description?: string) => {
      context.addToast({ type: "warning", title, description });
    },
    info: (title: string, description?: string) => {
      context.addToast({ type: "info", title, description });
    },
  };

  return { ...context, toast };
}

// Toast icons
const icons: Record<ToastType, React.ReactNode> = {
  success: (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
    </svg>
  ),
  error: (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
    </svg>
  ),
  warning: (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
    </svg>
  ),
  info: (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
  ),
};

// Toast styles
const styles: Record<ToastType, { container: string; icon: string }> = {
  success: {
    container: "bg-green-50 border-green-200 text-green-800",
    icon: "text-green-600",
  },
  error: {
    container: "bg-red-50 border-red-200 text-red-800",
    icon: "text-red-600",
  },
  warning: {
    container: "bg-amber-50 border-amber-200 text-amber-800",
    icon: "text-amber-600",
  },
  info: {
    container: "bg-sky-50 border-sky-200 text-sky-800",
    icon: "text-sky-600",
  },
};

// Toast component
function ToastItem({ toast, onRemove }: { toast: Toast; onRemove: () => void }) {
  const style = styles[toast.type];

  return (
    <div
      className={clsx(
        "flex items-start gap-3 p-4 rounded-lg border shadow-lg animate-slide-up",
        style.container
      )}
    >
      <span className={style.icon}>{icons[toast.type]}</span>
      <div className="flex-1 min-w-0">
        <p className="font-medium">{toast.title}</p>
        {toast.description && (
          <p className="mt-1 text-sm opacity-80">{toast.description}</p>
        )}
      </div>
      <button
        onClick={onRemove}
        className="shrink-0 p-1 rounded hover:bg-black/5 transition-colors"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  );
}

// Toast container
function ToastContainer() {
  const context = React.useContext(ToastContext);
  if (!context) return null;

  if (context.toasts.length === 0) return null;

  return (
    <div className="fixed bottom-4 right-4 z-50 flex flex-col gap-2 w-full max-w-sm">
      {context.toasts.map((toast) => (
        <ToastItem
          key={toast.id}
          toast={toast}
          onRemove={() => context.removeToast(toast.id)}
        />
      ))}
    </div>
  );
}

// Toast notification component (for standalone use)
interface ToastNotificationProps {
  type: ToastType;
  title: string;
  description?: string;
  onClose?: () => void;
}

export function ToastNotification({ type, title, description, onClose }: ToastNotificationProps) {
  const style = styles[type];

  return (
    <div
      className={clsx(
        "flex items-start gap-3 p-4 rounded-lg border shadow-lg",
        style.container
      )}
    >
      <span className={style.icon}>{icons[type]}</span>
      <div className="flex-1 min-w-0">
        <p className="font-medium">{title}</p>
        {description && <p className="mt-1 text-sm opacity-80">{description}</p>}
      </div>
      {onClose && (
        <button
          onClick={onClose}
          className="shrink-0 p-1 rounded hover:bg-black/5 transition-colors"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      )}
    </div>
  );
}
