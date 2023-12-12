import React from "react";
import MessageInput from "@/components/MessageInput";
import useWebSocket from "@/hook/useWebSocket";
import { getUserUUID } from "@/utils/liarHelper";
import ChatTimeline from "@/components/ChatTimeline";
import PlayAndWaitUserList from "@/components/PlayAndWaitUserList";
import ReadyButton from "@/components/ReadyButton";
import { useAtom } from "jotai";
import { isGameStartedAtom, isUserJoinGameAtom } from "@/store/gameAtom";
import VerticalLine from "@/components/Line/VerticalLine";

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
    <main className="max-w-screen-md mx-auto bg-white shadow-lg py-8 px-4 lg:px-8 min-h-screen">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-4xl font-bold italic">Liar of Turing</h1>
        {!isGameStarted && user && user.player_type !== "player" && (
          <ReadyButton
            userData={user}
            sendMessage={handleWebSocketMessageSend}
          />
        )}
      </div>
      {/* <hr className="w-[50%] h-1 bg-black" /> */}
      {/* <div className="flex flex-row-reverse justify-between"> */}
      <div className="flex flex-col lg:flex-row justify-between space-y-4 lg:space-y-0 lg:space-x-4">
        <ChatTimeline messageLogList={messageLogList} userData={user} />
        <VerticalLine />
        <PlayAndWaitUserList userData={user} />
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
