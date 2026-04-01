import { useEffect, useState } from "react";
import {
  createMessage,
  getCurrentUser,
  listChannels,
  listMessages,
  listSpaces,
  login,
} from "./lib/api";
import type { Message } from "./types";

type Space = {
  id: string;
  name: string;
  slug: string;
  description: string;
  visibility: string;
  discoverable: boolean;
  created_by: string;
  created_at: string;
  updated_at: string;
};

type Channel = {
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

export default function App() {
  const [email, setEmail] = useState("authjoe@example.com");
  const [password, setPassword] = useState("password123");
  const [token, setToken] = useState<string | null>(null);
  const [userName, setUserName] = useState<string>("none");

  const [spaces, setSpaces] = useState<Space[]>([]);
  const [selectedSpaceId, setSelectedSpaceId] = useState<string | null>(null);

  const [channels, setChannels] = useState<Channel[]>([]);
  const [selectedChannelId, setSelectedChannelId] = useState<string | null>(null);

  const [messages, setMessages] = useState<Message[]>([]);
  const [messageBody, setMessageBody] = useState("");
  const [nextCursor, setNextCursor] = useState<string | null>(null);

  const [error, setError] = useState<string | null>(null);

  async function handleLogin() {
    try {
      setError(null);

      const loginResult = await login(email, password);
      setToken(loginResult.token);

      const user = await getCurrentUser(loginResult.token);
      setUserName(user.display_name);

      const spaceRows = await listSpaces(loginResult.token);
      setSpaces(spaceRows);

      if (spaceRows.length > 0) {
        setSelectedSpaceId(spaceRows[0].id);
      } else {
        setSelectedSpaceId(null);
        setChannels([]);
        setSelectedChannelId(null);
        setMessages([]);
        setNextCursor(null);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "login failed");
    }
  }

  useEffect(() => {
    if (!token || !selectedSpaceId) {
      return;
    }

    async function loadChannels() {
      try {
        setError(null);
        const channelRows = await listChannels(selectedSpaceId, token);
        setChannels(channelRows);

        if (channelRows.length > 0) {
          setSelectedChannelId(channelRows[0].id);
        } else {
          setSelectedChannelId(null);
          setMessages([]);
          setNextCursor(null);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "failed to load channels");
      }
    }

    loadChannels();
  }, [selectedSpaceId, token]);

  useEffect(() => {
    if (!token || !selectedChannelId) {
      return;
    }

    async function loadChannelMessages() {
      try {
        setError(null);
        const result = await listMessages(selectedChannelId, token);
        setMessages(result.items);
        setNextCursor(result.next_cursor);
      } catch (err) {
        setError(err instanceof Error ? err.message : "failed to load messages");
      }
    }

    loadChannelMessages();
  }, [selectedChannelId, token]);

  async function handleLoadOlder() {
    if (!token || !selectedChannelId || !nextCursor) {
      return;
    }

    try {
      setError(null);
      const result = await listMessages(selectedChannelId, token, nextCursor);
      setMessages((previous) => [...result.items, ...previous]);
      setNextCursor(result.next_cursor);
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to load older messages");
    }
  }

  async function handleSendMessage() {
    if (!token || !selectedChannelId) {
      return;
    }

    const trimmed = messageBody.trim();
    if (!trimmed) {
      return;
    }

    try {
      setError(null);
      const created = await createMessage(selectedChannelId, token, trimmed);
      setMessages((previous) => {
        const alreadyExists = previous.some((message) => message.id === created.id);
        if (alreadyExists) {
          return previous;
        }
        return [...previous, created];
      });
      setMessageBody("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to send message");
    }
  }

  return (
    <div style={{ padding: "24px", fontFamily: "Arial, sans-serif" }}>
      <h1>manumental-effort</h1>

      <div style={{ display: "grid", gridTemplateColumns: "260px 260px 260px 1fr", gap: "16px" }}>
        <section>
          <h2>Login</h2>

          <label style={{ display: "block", marginBottom: "12px" }}>
            Email
            <input
              style={{ display: "block", width: "100%", marginTop: "4px" }}
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </label>

          <label style={{ display: "block", marginBottom: "12px" }}>
            Password
            <input
              type="password"
              style={{ display: "block", width: "100%", marginTop: "4px" }}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </label>

          <button onClick={handleLogin}>Login</button>

          <div style={{ marginTop: "16px" }}>
            <div>Token: {token ? "present" : "missing"}</div>
            <div>User: {userName}</div>
            <div>Selected space: {selectedSpaceId ?? "none"}</div>
            <div>Selected channel: {selectedChannelId ?? "none"}</div>
          </div>
        </section>

        <section>
          <h2>Spaces</h2>
          {spaces.length === 0 ? (
            <p>No spaces.</p>
          ) : (
            <ul>
              {spaces.map((space) => (
                <li key={space.id}>
                  <button onClick={() => setSelectedSpaceId(space.id)}>
                    {space.name} {space.id === selectedSpaceId ? "(selected)" : ""}
                  </button>
                </li>
              ))}
            </ul>
          )}
        </section>

        <section>
          <h2>Channels</h2>
          {channels.length === 0 ? (
            <p>No channels.</p>
          ) : (
            <ul>
              {channels.map((channel) => (
                <li key={channel.id}>
                  <button onClick={() => setSelectedChannelId(channel.id)}>
                    {channel.name} {channel.id === selectedChannelId ? "(selected)" : ""}
                  </button>
                </li>
              ))}
            </ul>
          )}
        </section>

        <section>
          <h2>Messages</h2>

          <div style={{ marginBottom: "12px" }}>
            <button onClick={handleLoadOlder} disabled={!nextCursor}>
              Load older
            </button>
          </div>

          <div
            style={{
              border: "1px solid #ccc",
              minHeight: "280px",
              maxHeight: "420px",
              overflow: "auto",
              padding: "12px",
              marginBottom: "12px",
            }}
          >
            {messages.length === 0 ? (
              <p>No messages.</p>
            ) : (
              messages.map((message) => (
                <div key={message.id} style={{ padding: "8px 0", borderBottom: "1px solid #eee" }}>
                  <div style={{ fontSize: "12px", marginBottom: "4px" }}>
                    {message.user_id} · {new Date(message.created_at).toLocaleString()}
                  </div>
                  <div>{message.body}</div>
                </div>
              ))
            )}
          </div>

          <div>
            <textarea
              rows={4}
              style={{ width: "100%", boxSizing: "border-box", marginBottom: "8px" }}
              value={messageBody}
              onChange={(e) => setMessageBody(e.target.value)}
              placeholder="Write a message"
            />
            <button onClick={handleSendMessage} disabled={!selectedChannelId}>
              Send
            </button>
          </div>
        </section>
      </div>

      {error ? <div style={{ marginTop: "16px", color: "red" }}>{error}</div> : null}
    </div>
  );
}
