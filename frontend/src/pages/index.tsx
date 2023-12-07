"use client";
import { useState, useEffect, use } from "react";
import { atomWithReset } from "jotai/utils";
import { useAtom } from "jotai";
import Image from "next/image";
import { Inter } from "next/font/google";
import { Message, Player } from "@/types/playerTypes";
import { randomUUID } from "crypto";
import { getUserUUID } from "@/utils/liarHelper";
import { WsJsonRequest, WsJsonResponse } from "@/types/wsTypes";
import { sendEnterHuman, sendLeftUser } from "@/utils/weHelper";
import * as datetime from "date-fns";

const getDayFromTimestamp = (timestamp: number): string => {
  return datetime.format(new Date(timestamp), "yyyy-MM-dd HH:mm:ss");
};

const inter = Inter({ subsets: ["latin"] });
const messageAtom = atomWithReset<string>("");
const chatLogAtom = atomWithReset<Message[]>([]);
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

  const setupWebSocket = () => {
    if (!socket) {
      const newSocket = new WebSocket(WEBSOCKET_URL);
      setSocket(newSocket);

      newSocket.onopen = () => {
        console.log("connected!!");
        sendEnterHuman(newSocket, userUUID);
      };

      newSocket.onmessage = handleWebSocketMessage;
      newSocket.onclose = handleWebSocketClose;

      setIsConnected(true);
    }
  };

  const handleWebSocketMessage = (event: MessageEvent) => {
    const data = JSON.parse(event.data);
    console.log("Received action:", data.action);
    switch (data.action) {
      case "human_info":
        console.log("human_info", data);
        console.log("username:", username, userData.username);
        console.log("userUUID:", userUUID);
        console.log("data.user.uuid:", data.user.uuid);
        setUserList(() => {
          if (
            userData.username.length !== 0 &&
            data.user.uuid !== userData.uuid
          )
            return [
              userData,
              ...data.online_user_list.filter((user: Player) => {
                return user.username !== userData.username;
              }),
            ];
          if (!data.user) return [data.online_user_list];
          return [
            data.user,
            ...data.online_user_list.filter((user: Player) => {
              return user.username !== data.user.username;
            }),
          ];
        });
        if (userUUID === data.user.uuid) {
          setUserData(data.user);
          setUsername(data.user.username);
        }
        setChatLog((chatLog: Message[]) => {
          const chat = {
            message_id: data.message_id,
            timestamp: data.timestamp,
            user: data.user,
            message: data.message,
            message_type: data.message_type,
          };
          return [...chatLog, chat];
        });

        break;
      case "user_list":
        console.log("user_list", data.online_user_list);

        setUserList(() => {
          if (!userData) return [data.online_user_list];
          return [
            userData,
            ...data.online_user_list.filter((user: Player) => {
              return user.username !== userData.username;
            }),
          ];
        });
        break;
      case "new_message":
        console.log("Message", data);
        setChatLog((chatLog: Message[]) => {
          const chat = {
            message_id: data.message_id,
            timestamp: data.timestamp,
            user: data.user,
            message: data.message,
            message_type: data.message_type,
          };
          return [...chatLog, chat];
        });

        break;
      case "update_state":
        console.log("update_state", data);
        setUserList(() => {
          if (!userData) return [data.online_user_list];
          return [
            userData,
            ...data.online_user_list.filter((user: Player) => {
              return user.username !== userData.username;
            }),
          ];
        });
        setChatLog((chatLog: Message[]) => {
          const chat = {
            message_id: data.message_id,
            timestamp: data.timestamp,
            user: data.user,
            message: data.message,
            message_type: data.message_type,
          };
          return [...chatLog, chat];
        });

        break;
    }
  };

  const handleWebSocketClose = () => {
    console.log("WebSocket closed");
    sendLeftUser(socket);
  };

  useEffect(() => {
    // setupWebSocket();
    // if (socket) return;
    const sc = new WebSocket(WEBSOCKET_URL);
    setSocket(sc);
    setIsConnected(true);
    return () => {
      // console.log("clean up");
      // sc?.close();
    };
  }, []);
  useEffect(() => {
    if (socket) {
      socket.onopen = () => {
        console.log("connected!!");
        sendEnterHuman(socket, userUUID);
      };
      socket.onmessage = handleWebSocketMessage;
      socket.onclose = handleWebSocketClose;
    }
  }, [socket, chatLog, userData, username]);

  useEffect(() => {
    const handleBeforeUnload = (event: BeforeUnloadEvent) => {
      console.log("Leaving");
      sendLeftUser(socket);
      event.preventDefault();
      event.returnValue = "";
    };
    window.addEventListener("beforeunload", handleBeforeUnload);
    return () => {
      console.log("clean up");
      window.removeEventListener("beforeunload", handleBeforeUnload);
    };
  }, [socket, userUUID]);

  // useEffect(() => {
  //   handleWebSocketMessage();
  // }, [chatLog, socket, userData, username]);

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
      <ul className="overflow-y-scroll flex-1 my-4 ">
        {chatLog.length > 0 &&
          chatLog.map((message: Message, idx) => {
            return (
              <li
                key={idx}
                className="flex py-0.5 pr-16 leading-[22px] hover:bg-gray-950/[.07]"
              >
                <div className="flex py-0.5 leading-[22px] hover:bg-gray-950/[.07]">
                  <div className="overflow-hidden relative mt-0.5 mr-2 w-10 min-w-fit h-10 rounded-full">
                    <Image
                      src="/favicon.png"
                      alt="discord"
                      layout="fill"
                      objectFit="contain"
                    />
                  </div>
                  {/* <span
                  className="font-bold mr-2"
                  style={{
                    color:
                      message.user.username === userData.username
                        ? "#3b82f6"
                        : "black",
                  }}
                >{`${message.user.username}>`}</span>
                <span>{message.message}</span> */}
                  <div>
                    <p className="flex items-baseline">
                      <span
                        className="mr-2 font-bold text-green-400"
                        style={{
                          color:
                            message.user.username === userData.username
                              ? "#3b82f6"
                              : "black",
                        }}
                      >
                        {message.user.username}
                      </span>
                      <span className="text-xs font-medium text-gray-900">
                        {getDayFromTimestamp(message.timestamp)}
                      </span>
                    </p>
                    <p className="text-gray-900">{message.message}</p>
                  </div>
                </div>
              </li>
            );
          })}
      </ul>
      <form className="mt-4 flex flex-row" onSubmit={handleSendMessage}>
        <label htmlFor="message">
          {userData.username && (
            <span
              className="mr-2 font-bold"
              // style={{
              //   color: "#3b82f6",
              // }}
            >{`${userData.username}: `}</span>
          )}
        </label>
        <div className="flex flex-row">
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
      {/* discord */}
    </main>
  );
}
