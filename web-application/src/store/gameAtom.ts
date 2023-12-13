import { atomWithReset } from "jotai/utils";
import { GAME_ROUND_NUM, GAME_TURNS_NUM, MAX_PLAYER } from "./gameStore";
import { playerListAtom, userAtom } from "./chatAtom";
import { atom, useAtomValue } from "jotai";
import { UserSelection } from "@/types/wsTypes";
import { initialUserSelection } from "./chatStore";

export const maxPlayerAtom = atom<number>(MAX_PLAYER);

// export const isGameStartedAtom = atom<boolean>((get) => {
//   const playerList = get(playerListAtom);
//   const maxPlayer = get(maxPlayerAtom);
//   return playerList.length === maxPlayer;
// });
export const isGameStartedAtom = atomWithReset<boolean>(false);

export const isUserJoinGameAtom = atom<boolean>((get) => {
  const isGameStarted = get(isGameStartedAtom);
  const user = get(userAtom);
  return isGameStarted && user.player_type === "player";
});

export const isYourTurnAtom = atom<boolean>(false);
export const gameRoundAtom = atom<number>(0);
export const gameRoundNumAtom = atom<number>(GAME_ROUND_NUM);
export const gameTurnsLeftAtom = atom<number>(0);
export const gameTurnsNumAtom = atom<number>(GAME_TURNS_NUM);

export const isFinishedRoundAtom = atom<boolean>(false);

export const isFinishedSubmitionAtom = atom<boolean>(false);

export const userSelectionAtom = atom<UserSelection>(initialUserSelection);
export const userSelectionListAtom = atom<UserSelection[]>([]);
