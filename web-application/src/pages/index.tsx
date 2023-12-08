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
import { time } from "console";
import { TIMER_TIME } from "@/store/gameStore";

const getDayFromTimestamp = (timestamp: number): string => {
  return datetime.format(new Date(timestamp), "yyyy-MM-dd HH:mm:ss");
};
const getTimeFromTimer = (time: number): string => {
  //using date-fns ans 10:00
  const date = new Date(time * 1000);
  let minutes = (date.getMinutes() < 10 ? "0" : "") + date.getMinutes();
  let seconds = (date.getSeconds() < 10 ? "0" : "") + date.getSeconds();
  return `${minutes}:${seconds}`;
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
const timerTimeAtom = atomWithReset<number>(0);
const timerIsRunningAtom = atomWithReset<boolean>(false);
const timerIsPausedAtom = atomWithReset<boolean>(false);
const isGameStartedAtom = atomWithReset<boolean>(false);
let userUUID: string = getUserUUID();
const WEBSOCKET_URL = process.env.NEXT_PUBLIC_WEBSOCKET_URL
  ? process.env.NEXT_PUBLIC_WEBSOCKET_URL
  : "";

const isDebugModeAtom = atomWithReset<boolean>(false);
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

  const [timerTime, setTimerTime] = useAtom(timerTimeAtom);
  const [timerIsRunning, setTimerIsRunning] = useAtom(timerIsRunningAtom);
  const [timerIsPaused, setTimerIsPaused] = useAtom(timerIsPausedAtom);
  const [isGameStarted, setIsGameStarted] = useAtom(isGameStartedAtom);

  // for test
  const [testUsername, setTestUsername] = useState<string>("");
  const [testMessage, setTestMessage] = useState<string>("");
  const [isDebugMode, setIsDebugMode] = useAtom(isDebugModeAtom);

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

  const handleChangeTestMessage = (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    setTestMessage(event.target.value);
  };

  const handleTestSendMessage = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const message = event.currentTarget.message.value.trim();
    if (message.length === 0) return;

    const sendUser = userList.find((user: Player) => {
      return user.username === testUsername;
    });
    if (!sendUser) return;
    const jsonData = {
      action: "new_message",
      room_id: 0,
      user: sendUser,
      timestamp: Date.now(),
      message: message,
    };
    socket?.send(JSON.stringify(jsonData));
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
          if (!data.online_user_list) return [];

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
          if (!data.online_user_list) return [];
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
          if (!data.online_user_list) return [];
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

  useEffect(() => {
    const container = window.document.getElementById("chatLog");
    if (!container) return;
    container.scrollTop = container.scrollHeight;
  }, [chatLog]);

  // useEffect(() => {
  //   handleWebSocketMessage();
  // }, [chatLog, socket, userData, username]);

  useEffect(() => {
    if (timerTime === 0) {
      setTimerIsRunning(false);
      setTimerIsPaused(false);
    }
  }, [timerTime]);
  useEffect(() => {
    if (!timerIsRunning) return;
    const interval = setInterval(() => {
      setTimerTime((timerTime) => timerTime - 1);
    }, 1000);
    return () => clearInterval(interval);
  }, [timerIsRunning]);

  if (!isConnected) return <div>Connecting...</div>;
  return (
    <main
      className={`py-8 mx-auto ${inter.className} w-[80vw] max-w-2xl min-h-max relative`}
    >
      <div className="flex flex-row justify-between">
        <h1 className="text-3xl font-bold">Liar of Turing</h1>
        <div className="flex flex-row">
          <div className="text-center text-xl font-medium mx-2 text-gray-900">
            {getTimeFromTimer(timerTime)}
          </div>
          <button
            className="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md"
            onClick={() => {
              setTimerTime(TIMER_TIME);
              setTimerIsRunning(true);
              if (!isGameStarted) setIsGameStarted(true);
            }}
          >
            Start
          </button>
        </div>
      </div>
      <div className="flex flex-row-reverse justify-between">
        <div className="flex flex-col ml-4">
          <h3 className="mt-6 font-bold text-xl">Now online</h3>
          <ul className="my-4 h-80 flex-1">
            {userList &&
              userList.length > 0 &&
              userList.map((user: Player, index) => {
                return (
                  <li
                    key={index}
                    style={{
                      color:
                        user.username === userData.username
                          ? "#3b82f6"
                          : "black",
                      fontWeight:
                        user.username === userData.username ? 700 : 400,
                    }}
                  >
                    {user.username}
                  </li>
                );
              })}
          </ul>
        </div>
        <div className="flex flex-col flex-1">
          <h3 className="mt-6 font-bold text-xl">Chat Log</h3>
          <ul
            className="overflow-y-scroll min-h-[60vh] max-h-[60vh] my-4 flex-1"
            id="chatLog"
          >
            {chatLog.length > 0 &&
              chatLog.map((message: Message, idx) => {
                if (message.message_type === "system") return;
                return (
                  <li key={idx} className="flex py-0.5 pr-16 leading-[22px]">
                    <div className="flex py-0.5 leading-[22px]">
                      <div className="overflow-hidden relative mt-0.5 mr-2 w-10 min-w-fit h-10 rounded-sm">
                        <Image
                          src={`/nickname_icon/${message.user.nickname_id}.png`}
                          alt={`${message.user.username} icon`}
                          layout="fill"
                          objectFit="contain"
                        />
                      </div>
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
        </div>
      </div>
      {userData.username && (
        <form className="mt-4 flex flex-row" onSubmit={handleSendMessage}>
          <label htmlFor="message">
            {userData.username && (
              <span
                className="mr-2 font-bold flex-1"
                style={{
                  color: "#3b82f6",
                }}
              >{`${userData.username}: `}</span>
            )}
          </label>
          <input
            autoFocus
            className="border-2 border-gray-400 rounded-md flex-1"
            type="text"
            id="message"
            value={nowMessage}
            onChange={handleChangeMessage}
          />
        </form>
      )}
      {isGameStarted && timerTime === 0 && (
        <div className="flex flex-col">
          <h3 className="mt-6 font-bold text-xl">Choose AI</h3>
          <form className="flex flex-row" onSubmit={handleTestSendMessage}>
            <label htmlFor="ai">{`I think the AI is..`}</label>
            <select
              name="ai"
              id="ai"
              onChange={(e) => {
                setTestUsername(e.target.value);
              }}
              className="border-2 border-gray-400 rounded-md w-fit-content"
            >
              {userList &&
                userList.length > 0 &&
                userList.map((user: Player, index) => {
                  return (
                    <option key={index} value={user.username}>
                      {user.username}
                    </option>
                  );
                })}
            </select>
            <label htmlFor="reason" className="mx-2 ">
              Reason
            </label>
            <input
              className="border-2 border-gray-400 rounded-md w-fit-content"
              type="text"
              id="reason"
              value={testMessage}
              onChange={handleChangeTestMessage}
            />
          </form>
        </div>
      )}

      <button
        className="top-0 -right-[25vw] absolute text-white hover:text-black"
        onClick={() => setIsDebugMode((isDebugMode) => !isDebugMode)}
      >
        debug
      </button>

      {/* for test */}
      {isDebugMode && (
        <div className="flex flex-col">
          <div className="mt-40"></div>

          <form className="flex flex-row" onSubmit={handleTestSendMessage}>
            <label htmlFor="username">Username</label>
            <select
              name="username"
              id="username"
              onChange={(e) => {
                setTestUsername(e.target.value);
              }}
              className="border-2 border-gray-400 rounded-md w-fit-content"
            >
              {userList &&
                userList.length > 0 &&
                userList.map((user: Player, index) => {
                  return (
                    <option key={index} value={user.username}>
                      {user.username}
                    </option>
                  );
                })}
            </select>
            <label htmlFor="message" className="mx-2 ">
              Message
            </label>
            <input
              className="border-2 border-gray-400 rounded-md w-fit-content"
              type="text"
              id="message"
              value={testMessage}
              onChange={handleChangeTestMessage}
            />
          </form>
        </div>
      )}
    </main>
  );
}
