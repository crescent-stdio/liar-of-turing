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
  isFinishedSubmitionAtom,
  isGameStartedAtom,
  isUserJoinGameAtom,
  isYourTurnAtom,
} from "@/store/gameAtom";
import VerticalLine from "@/components/Line/VerticalLine";
import HorizontalLine from "@/components/Line/HorizontalLine";
import ShowGameStatus from "@/components/ShowGameStatus";
import ChooseAIInput from "@/components/game/ChooseAIInput";
import WaitingForSelection from "@/components/game/WaitingForSelection";
import { playerListAtom, userAtom } from "@/store/chatAtom";
import FinishedRoundModal from "@/components/game/FinishedRoundModal";
import InputComponent from "@/components/game/InputModal";
import InputModal from "@/components/game/InputModal";

export default function Page() {
  const { isConnected, messageLogList, handleWebSocketMessageSend } =
    useWebSocket(getUserUUID(), null);
  const [user, setUser] = useAtom(userAtom);
  const [isGameStarted] = useAtom(isGameStartedAtom);
  const [isYourTurn] = useAtom(isYourTurnAtom);

  if (!isConnected) return <div>Connecting...</div>;
  console.log(user);
  return (
    <main className="max-w-screen-md mx-auto bg-white shadow-lg py-4 md:py-8 px-4 md:px-8 min-h-screen">
      <div className="max-full min-h-max relative">
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
        <div className="flex flex-col-reverse md:flex-row justify-between space-y-4 md:space-y-0 md:space-x-4 mt-4">
          <ChatTimeline messageLogList={messageLogList} userData={user} />
          {/* <VerticalLine /> */}
          <PlayAndWaitUserList userData={user} />
        </div>
        <InputModal
          isGameStarted={isGameStarted}
          isYourTurn={isYourTurn}
          sendMessage={handleWebSocketMessageSend}
        />

        <FinishedRoundModal sendMessage={handleWebSocketMessageSend} />
      </div>
    </main>
  );
}
