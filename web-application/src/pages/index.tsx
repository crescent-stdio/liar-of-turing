import React, { useEffect } from "react";
import useWebSocket from "@/hook/useWebSocket";
import { getUserUUID } from "@/utils/liarHelper";
import ChatTimeline from "@/components/ChatTimeline";
import PlayAndWaitUserList from "@/components/PlayAndWaitUserList";
import ReadyButton from "@/components/ReadyButton";
import { useAtom } from "jotai";
import {
  isFinishedShowResultAtom,
  isGameStartedAtom,
  isYourTurnAtom,
} from "@/store/gameAtom";
import ShowGameStatus from "@/components/ShowGameStatus";
import { messageLogListAtom, userAtom } from "@/store/chatAtom";
import FinishedRoundModal from "@/components/game/FinishedRoundModal";
import InputModal from "@/components/game/InputModal";
import { RESULT_OPEN_TIME } from "@/store/gameStore";
import { WsJsonRequest } from "@/types/wsTypes";
import { initialUserSelection } from "@/store/chatStore";

export default function Page() {
  const { isConnected, messageLogList, handleWebSocketMessageSend } =
    useWebSocket(getUserUUID(), null);
  const [user, setUser] = useAtom(userAtom);
  const [isGameStarted] = useAtom(isGameStartedAtom);
  const [isYourTurn] = useAtom(isYourTurnAtom);
  const [isFinishedShowResult, setIsFinishedShowResult] = useAtom(
    isFinishedShowResultAtom
  );
  const [, setMessageLogList] = useAtom(messageLogListAtom);

  useEffect(() => {
    if (!isFinishedShowResult) return;
    console.log("isFinishedShowResult", isFinishedShowResult);
    const timer = setTimeout(() => {
      const jsonData: WsJsonRequest = {
        action: "restart_game",
        user: user,
        timestamp: Date.now(),
        max_player: 0,
        message: "",
        game_round: 0,
        game_turns_left: 0,
        game_round_num: 0,
        game_turn_num: 0,
        user_selection: initialUserSelection,
      };
      handleWebSocketMessageSend(jsonData);
      setIsFinishedShowResult(false);
    }, RESULT_OPEN_TIME);
    return () => clearTimeout(timer);
  }, [isFinishedShowResult]);

  if (!isConnected) return <div>Connecting...</div>;
  // console.log(user);
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
          <PlayAndWaitUserList />
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
