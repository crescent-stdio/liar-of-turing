"use client";
import { useState, useEffect, use } from "react";
import { atomWithReset } from "jotai/utils";
import { useAtom } from "jotai";
import Image from "next/image";
import { Inter } from "next/font/google";
import { Chat, Message, Player } from "@/types/playerTypes";
import { randomUUID } from "crypto";
import { getUserUUID } from "@/utils/liarHelper";

const inter = Inter({ subsets: ["latin"] });

const messageAtom = atomWithReset<Chat>({
  messages: [],
});

const chatLogAtom = atomWithReset<Chat>({
  messages: [],
});
const userListAtom = atomWithReset<Player[]>([]);
const usernameAtom = atomWithReset<string>("");
const userDataAtom = atomWithReset<Player>({
  uuid: "",
  user_id: -1,
  username: "",
  role: "",
});

let userUUID: string = getUserUUID();
const WEBSOCKET_URL = process.env.NEXT_PUBLIC_WEBSOCKET_URL
  ? process.env.NEXT_PUBLIC_WEBSOCKET_URL
  : "";
export default function Home() {
  const [isConnected, setIsConnected] = useState<boolean | null>(null);
  const [count, setCount] = useState<number>(0);
  const [res, setRes] = useState<string>("");
  const [message, setMessage] = useAtom(messageAtom);
  const [nowMessage, setNowMessage] = useState<string>("");
  const [chatLog, setChatLog] = useAtom(chatLogAtom);
  const [userData, setUserData] = useAtom(userDataAtom);
  const [username, setUsername] = useAtom(usernameAtom);
  const [inputUsername, setInputUsername] = useState<string>("");
  const [userList, setUserList] = useAtom(userListAtom);
  const [socket, setSocket] = useState<WebSocket | null>(null);

  const handleChangeMessage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setNowMessage(event.target.value);
  };
  const handleSendMessage = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const message = event.currentTarget.message.value.trim();
    if (message.length === 0) return;

    const jsonData = {
      action: "new_message",
      room_id: 0,
      user: userData,
      timestamp: Date.now(),
      message: message,
    };
    socket?.send(JSON.stringify(jsonData));
    setNowMessage("");
  };
  useEffect(() => {
    // if (WEBSOCKET_URL.length === 0) return;
    // setUserList([]);
    if (!socket) {
      const sc = new WebSocket(WEBSOCKET_URL);
      setSocket(sc);

      setIsConnected(true);
      console.log(sc);
      return;
    }
    socket.onopen = () => {
      console.log("connected!!");

      const jsonData = {
        action: "enter_human",
        room_id: 0,
        user: {
          uuid: userUUID,
          user_id: -1,
          username: "",
          role: "human",
          is_online: true,
        },
        timestamp: Date.now(),
        message: "",
      };
      socket.send(JSON.stringify(jsonData));
    };
    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log("Action:", data.action);
      switch (data.action) {
        case "human_info":
          console.log("human_info", data);
          console.log("username:", username, userData.username);
          console.log("userUUID:", userUUID);
          console.log("data.user.uuid:", data.user.uuid);
          setUserList(data.connected_users);
          if (userUUID === data.user.uuid) {
            setUserData(data.user);
            setUsername(data.user.username);

            // send welcome message
          }
          setChatLog((chatLog: Chat) => {
            return {
              messages: [
                ...chatLog.messages,
                {
                  id: data.id,
                  timestamp: data.timestamp,
                  sender: data.user.username,
                  content: data.message,
                },
              ],
            };
          });

          break;
        case "user_list":
          console.log("user_list", data.connected_users);
          setUserList(data.connected_users);
          break;
        case "message":
          console.log("Message", data);
          setChatLog((chatLog: Chat) => {
            return {
              messages: [
                ...chatLog.messages,
                {
                  id: data.id,
                  timestamp: data.timestamp,
                  sender: data.user.username,
                  content: data.message,
                },
              ],
            };
          });
          break;
        case "update_state":
          console.log("update_state", data);
          setUserList(data.connected_users);
          setChatLog((chatLog: Chat) => {
            return {
              messages: [
                ...chatLog.messages,
                {
                  id: data.id,
                  timestamp: data.timestamp,
                  sender: data.user.username,
                  content: data.message,
                },
              ],
            };
          });
          break;
          break;
      }
    };
    socket.onclose = () => {
      const jsonData = { action: "left_user" };
      socket?.send(JSON.stringify(jsonData));
    };
    return () => {
      // const jsonData = { action: "left_user" };
      // socket?.send(JSON.stringify(jsonData));
      // console.log("clean up");
      // socket?.close();
    };
  }, [chatLog, socket, userData, username]);

  useEffect(() => {
    if (!socket) return;

    const handleBeforeUnload = (event: BeforeUnloadEvent) => {
      console.log("Leaving");
      const jsonData = { action: "left_user" };
      socket?.send(JSON.stringify(jsonData));
      // // Make sure to check if the socket is connected before trying to send
      // if (socket && socket.readyState === WebSocket.OPEN) {
      //   socket.send(JSON.stringify(jsonData));
      event.preventDefault();
      event.returnValue = "";
    };
    window.addEventListener("beforeunload", handleBeforeUnload);
    return () => {
      const jsonData = { action: "left_user" };
      socket?.send(JSON.stringify(jsonData));
      console.log("clean up");
      socket?.close();
      window.removeEventListener("beforeunload", handleBeforeUnload);
    };
  }, [socket, userUUID]);

  if (!isConnected) return <div>Connecting...</div>;
  return (
    <main className={`m-4 ${inter.className}`}>
      <h1 className="text-3xl font-bold">Liars of Turing</h1>
      <div className="flex flex-col">
        <h3 className="mb-2 font-bold text-xl">Now connected users</h3>
        {userList &&
          userList.length > 0 &&
          userList.map((user: Player, index) => {
            return (
              <p
                key={index}
                style={{
                  color:
                    user.username === userData.username ? "#3b82f6" : "black",
                  fontWeight: user.username === userData.username ? 700 : 400,
                }}
              >
                {user.username}
              </p>
            );
          })}
      </div>
      {/* {userData.username && (
          <h3 className="my-4 font-bold text-xl">
            {`Now your name: `}
            <span className="text-blue-500">{userData.username}</span>
          </h3>
        )} */}

      <h3 className="my-4 font-bold text-xl">Chat Log</h3>
      <ul className="w-96 max-w-1/2 h-80 border-2 border-gray-400 rounded-md overflow-y-scroll">
        {chatLog.messages.length > 0 &&
          chatLog.messages.map((message: Message, idx) => {
            return (
              <li key={idx} className="">
                <span
                  className="font-bold mr-2"
                  style={{
                    color:
                      message.sender === userData.username
                        ? "#3b82f6"
                        : "black",
                  }}
                >{`${message.sender}>`}</span>
                <span>{message.content}</span>
              </li>
            );
          })}
      </ul>
      <form className="flex flex-col" onSubmit={handleSendMessage}>
        <label htmlFor="message">Message</label>
        <div className="flex flex-row">
          {userData.username && (
            <span className="mr-2 font-bold">{`${userData.username}: `}</span>
          )}
          <input
            autoFocus
            className="border-2 border-gray-400 rounded-md w-fit-content"
            type="text"
            id="message"
            value={nowMessage}
            onChange={handleChangeMessage}
          />
        </div>
      </form>
    </main>
  );
}
