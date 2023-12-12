import React from "react";
import MessageInput from "@/components/MessageInput";
import useWebSocket from "@/hook/useWebSocket";
import { getUserUUID } from "@/utils/liarHelper";
import ChatTimeline from "@/components/ChatTimeline";
import PlayAndWaitUserList from "@/components/PlayAndWaitUserList";
import ReadyButton from "@/components/ReadyButton";
import { useAtom, useAtomValue } from "jotai";
import {
  isFinishedRoundAtom,
  isGameStartedAtom,
  isUserJoinGameAtom,
  isYourTurnAtom,
} from "@/store/gameAtom";
import VerticalLine from "@/components/Line/VerticalLine";
import HorizontalLine from "@/components/Line/HorizontalLine";
import ShowGameStatus from "@/components/ShowGameStatus";
import ChooseAIInput from "@/components/game/ChooseAIInput";
import WaitingForSelection from "@/components/game/WaitingForSelection";

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
  const isYourTurn = useAtomValue(isYourTurnAtom);
  const isFinishedRound = useAtomValue(isFinishedRoundAtom);
  if (!isConnected) return <div>Connecting...</div>;

  return (
    <main className="max-w-screen-md mx-auto bg-white shadow-lg py-4 md:py-8 px-4 md:px-8 min-h-screen">
      <div className="flex justify-between items-center mb:2 md:mb-6">
        <h1 className="text-4xl font-bold italic underline">{`Liar of Turing`}</h1>
        {!isGameStarted && user && user.player_type !== "player" && (
          <ReadyButton
            userData={user}
            sendMessage={handleWebSocketMessageSend}
          />
        )}
        {isGameStarted && <ShowGameStatus />}
      </div>
      {/* <hr className="w-[50%] h-1 bg-black" /> */}
      {/* <div className="flex flex-row-reverse justify-between"> */}
      <div className="flex flex-col-reverse md:flex-row justify-between space-y-4 md:space-y-0 md:space-x-4 mt-4">
        <ChatTimeline messageLogList={messageLogList} userData={user} />
        {/* <VerticalLine /> */}
        <PlayAndWaitUserList userData={user} />
      </div>
      {isYourTurn ? (
        <MessageInput
          userData={user}
          sendMessage={handleWebSocketMessageSend}
        />
      ) : (
        <HorizontalLine />
      )}
      {isFinishedRound && user && user.player_type === "player" && (
        <ChooseAIInput
          userData={user}
          userList={userList}
          sendMessage={handleWebSocketMessageSend}
        />
      )}
      {isFinishedRound && user && user.player_type !== "player" && (
        <WaitingForSelection />
      )}
    </main>
  );
}
