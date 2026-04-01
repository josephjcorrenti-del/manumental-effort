export type User = {
  id: string;
  username: string;
  display_name: string;
  email: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

export type Channel = {
  id: string;
  space_id: string;
  name: string;
  slug: string;
  description: string;
  visibility: string;
  created_by: string;
  created_at: string;
  updated_at: string;
};

export type Message = {
  id: string;
  channel_id: string;
  user_id: string;
  body: string;
  created_at: string;
  updated_at: string;
};

export type MessageListResponse = {
  items: Message[];
  next_cursor: string | null;
};

export type LoginResponse = {
  token: string;
};

export type RealtimeServerEvent = {
  type: string;
  channel_id?: string;
  message?: Message;
  error?: string;
};
