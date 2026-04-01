import type { Channel, LoginResponse, Message, MessageListResponse, User } from "../types";

const API_BASE_URL = "http://localhost:8081";

async function fetchJson<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(options?.headers ?? {}),
    },
  });

  const text = await response.text();
  const data = text ? JSON.parse(text) : null;

  if (!response.ok) {
    const message = data?.error ?? `request failed: ${response.status}`;
    throw new Error(message);
  }

  return data as T;
}

export async function login(email: string, password: string): Promise<LoginResponse> {
  return fetchJson<LoginResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

export async function getCurrentUser(token: string): Promise<User> {
  return fetchJson<User>("/auth/me", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
}

export async function listChannels(spaceId: string, token: string): Promise<Channel[]> {
  return fetchJson<Channel[]>(`/spaces/${spaceId}/channels`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
}

export async function listMessages(
  channelId: string,
  token: string,
  before?: string,
  limit = 50,
): Promise<MessageListResponse> {
  const params = new URLSearchParams();
  params.set("limit", String(limit));
  if (before) {
    params.set("before", before);
  }

  return fetchJson<MessageListResponse>(`/channels/${channelId}/messages?${params.toString()}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
}

export async function createMessage(
  channelId: string,
  token: string,
  body: string,
): Promise<Message> {
  return fetchJson<Message>(`/channels/${channelId}/messages`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ body }),
  });
}
