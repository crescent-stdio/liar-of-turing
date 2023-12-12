import React from "react";
import OnlineUserList from "@/components/OnlineUserList";
import MessageInput from "@/components/MessageInput";
import TimerDisplay from "@/components/TimeDisplay";
import DebugPanel from "@/components/settings/AdminMessageSender";
import useWebSocket from "@/hook/useWebSocket";
import useTimer from "@/hook/useTimer";
import { getUserUUID } from "@/utils/liarHelper";
import ChatTimeline from "@/components/ChatTimeline";
import PlayAndWaitUserList from "@/components/PlayAndWaitUserList";
import ReadyButton from "@/components/ReadyButton";
import { useAtom } from "jotai";
import { isGameStartedAtom, isUserJoinGameAtom } from "@/store/gameAtom";

export default function Page() {
  const {
    socket,
    isConnected,
    userList,
    user,
    messageLogList,
    handleWebSocketMessageSend,
  } = useWebSocket(getUserUUID(), null);
  const [isGameStarted] = useAtom(isGameStartedAtom);
  const [isUserJoinGame] = useAtom(isUserJoinGameAtom);

  if (!isConnected) return <div>Connecting...</div>;

  return (
    <main className="py-8 mx-auto w-[80vw] max-w-2xl min-h-max relative">
      <div className="flex flex-row justify-between">
        <h1 className="text-3xl font-bold italic underline">Liar of Turing</h1>
        {!isGameStarted && user && user.player_type !== "player" && (
          <ReadyButton
            userData={user}
            sendMessage={handleWebSocketMessageSend}
          />
        )}
      </div>
      {/* <hr className="w-[50%] h-1 bg-black" /> */}
      <div className="flex flex-row-reverse justify-between">
        <PlayAndWaitUserList userData={user} />
        <ChatTimeline messageLogList={messageLogList} userData={user} />
      </div>
      {isUserJoinGame && (
        <MessageInput
          userData={user}
          sendMessage={handleWebSocketMessageSend}
        />
      )}
    </main>
  );
}
