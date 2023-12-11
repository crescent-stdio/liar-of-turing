import React from "react";
import OnlineUserList from "@/components/OnlineUserList";
import MessageInput from "@/components/MessageInput";
import TimerDisplay from "@/components/TimeDisplay";
import DebugPanel from "@/components/DebugPanel";
import useWebSocket from "@/hook/useWebSocket";
import useTimer from "@/hook/useTimer";
import { getUserUUID } from "@/utils/liarHelper";
import ChatTimeline from "@/components/ChatTimeline";

export default function Page() {
  const {
    socket,
    isConnected,
    userList,
    user,
    messageLogList,
    handleWebSocketMessageSend,
  } = useWebSocket(getUserUUID());
  const { timerTime, startTimer } = useTimer();

  if (!isConnected) return <div>Connecting...</div>;

  return (
    <main className="py-8 mx-auto w-[80vw] max-w-2xl min-h-max relative">
      <div className="flex flex-row justify-between">
        <h1 className="text-3xl font-bold italic">Liar of Turing</h1>
        <TimerDisplay timerTime={timerTime} startTimer={startTimer} />
      </div>
      <div className="flex flex-row-reverse justify-between">
        <OnlineUserList userList={userList} userData={user} />
        <ChatTimeline messageLogList={messageLogList} userData={user} />
      </div>
      <MessageInput
        userData={user}
        onSendMessage={handleWebSocketMessageSend}
      />
      <DebugPanel
        userList={userList}
        sendMessage={handleWebSocketMessageSend}
      />
    </main>
  );
}
