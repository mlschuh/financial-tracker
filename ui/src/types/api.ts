// src/types/api.ts
export interface Account {
  id: string;
  name: string;
  color: string;
}

export interface AccountCreate {
  name: string;
  color: string;
}

export interface Exception {
  type: "single" | "forever" | "skip";
  amount?: number;
}

export interface Event {
  id: string;
  name: string;
  category: string;
  account: string;
  amount: number;
  start: string;
  rrule?: string;
  type: "income" | "expense";
  exceptions?: { [key: string]: Exception };
}

export interface EventCreate {
  name: string;
  category: string;
  account: string;
  amount: number;
  start: string;
  rrule?: string;
  type: "income" | "expense";
  exceptions?: { [key: string]: Exception };
}

export interface EventOccurrence {
  id: string;
  date: string;
  amount: number;
  eventId: string;
  accountId: string;
  eventType: "income" | "expense";
  eventName: string;
}

export interface AccountBalance {
  date: string;
  balance: number;
  accountId: string;
  eventId: string;
}

export interface AppStateResponse {
  eventOccurances: EventOccurrence[];
  accountBalances: AccountBalance[];
  events: Event[];
  accounts: Account[];
}

export interface ErrorResponse {
  error: string;
}

export type ToastType = "success" | "error" | "warning" | "info";
