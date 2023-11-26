"use client";
import { useState, useEffect } from "react";
import { atomWithReset } from "jotai/utils";
import { useAtom } from "jotai";

interface Counter {
  value: number;
}

interface Player {
  id: number;
  name: string;
  isLiar: boolean;
}

interface Game {
  players: Player[];
  isOver: boolean;
}

interface Message {
  id: number;
  timestamp: number;
  sender: string;
  content: string;
}
interface Chat {
  messages: Message[];
}

const messageAtom = atomWithReset<Chat>({
  messages: [],
});

const chatLogAtom = atomWithReset<Chat>({
  messages: [],
});
const usernameAtom = atomWithReset<string>("");
export default function Home() {
  const [count, setCount] = useState<number>(0);
  const [res, setRes] = useState<string>("");
  const [message, setMessage] = useAtom(messageAtom);
  const [chatLog, setChatLog] = useAtom(chatLogAtom);
  const [username, setUsername] = useAtom(usernameAtom);
  const [inputUsername, setInputUsername] = useState<string>("");
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const handleUsernameSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setUsername(inputUsername);
  };

  const handleUsernameInput = (event: React.ChangeEvent<HTMLInputElement>) => {
    setInputUsername(event.target.value);
  };

  const connectTestHandler = async () => {
    if (!socket) return;
    socket.onmessage = (event) => {
      console.log(event.data);
      const data = JSON.parse(event.data);
      setChatLog((chatLog) => {
        return {
          messages: [
            ...chatLog.messages,
            {
              id: 0,
              timestamp: Date.now(),
              sender: "server",
              content: data.message,
            },
          ],
        };
      });
    };
  };
  useEffect(() => {
    try {
      const sc = new WebSocket("ws://localhost:8080/ws");
      sc.onopen = () => {
        console.log("connected");
      };
      setSocket(sc);
      console.log(sc);
    } catch (error) {
      console.error("Failed to fetch counter value:", error);
    }
  }, []);

  if (socket) {
    socket.onmessage = (event) => {
      console.log(event.data);
      const data = JSON.parse(event.data);
      setChatLog((chatLog) => {
        return {
          messages: [
            ...chatLog.messages,
            {
              id: 0,
              timestamp: Date.now(),
              sender: "server",
              content: data.message,
            },
          ],
        };
      });
    };
    socket.onclose = (event) => {
      console.log("closed");
    };
    socket.onerror = (event) => {
      console.log("error");
    };
    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log(data);
      setChatLog((chatLog) => {
        return {
          messages: [
            ...chatLog.messages,
            {
              id: 0,
              timestamp: Date.now(),
              sender: "server",
              content: data.message,
            },
          ],
        };
      });
    };
  }

  return (
    <main className="m-4">
      <h1 className="text-3xl font-bold">Liars of Turing</h1>
      <button
        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        onClick={connectTestHandler}
      >
        Connect to server
      </button>
      {res}
      <form className="flex flex-col" onSubmit={handleUsernameSubmit}>
        <label htmlFor="name">Name</label>
        <input
          type="text"
          id="name"
          value={inputUsername}
          onChange={handleUsernameInput}
        />
      </form>
      <ul>
        {chatLog.messages.map((message) => {
          return (
            <li key={message.id}>
              <p className="font-bold">{message.sender}</p>
              <p>{message.content}</p>
            </li>
          );
        })}
      </ul>
      <form className="flex flex-col">
        <label htmlFor="message">Message</label>
        <input type="text" id="message" />
      </form>
    </main>
  );
}
