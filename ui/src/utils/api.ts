// src/utils/api.ts
import axios from "axios";
import type {
  Account,
  AccountCreate,
  Event,
  EventCreate,
  AppStateResponse,
} from "@/types/api";

const api = axios.create({
  baseURL: "http://localhost:8080",
});

export const apiClient = {
  async getAppState(): Promise<AppStateResponse> {
    const response = await api.get<AppStateResponse>("/api/state");
    return response.data;
  },

  async getAccounts(): Promise<Account[]> {
    const response = await api.get<Account[]>("/api/accounts");
    return response.data;
  },

  async createAccount(account: AccountCreate): Promise<Account> {
    const response = await api.post<Account>("/api/accounts", account);
    return response.data;
  },

  async getEvents(): Promise<Event[]> {
    const response = await api.get<Event[]>("/api/events");
    return response.data;
  },

  async createEvent(event: EventCreate): Promise<Event> {
    const response = await api.post<Event>("/api/events", event);
    return response.data;
  },

  async deleteEvent(eventId: string): Promise<void> {
    await api.delete(`/api/events/${eventId}`);
  },
};
