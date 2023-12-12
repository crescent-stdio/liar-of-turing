import { User } from "@/types/playerTypes";
import OnlineUserList from "./OnlineUserList";
import { playerListAtom, watcherListAtom } from "@/store/chatAtom";
import { useAtomValue } from "jotai";
import { maxPlayerAtom } from "@/store/gameAtom";
import { useEffect, useState } from "react";
import HorizontalLine from "./Line/HorizontalLine";
type PlayAndWaitUserListProps = {
  userData: User;
};
export default function PlayAndWaitUserList({
  userData,
}: PlayAndWaitUserListProps) {
  const playerList = useAtomValue(playerListAtom);
  const watcherList = useAtomValue(watcherListAtom);
  const maxPlayer = useAtomValue(maxPlayerAtom);
  const [playerListTitle, setPlayerListTitle] = useState<string>("Players");

  // const playerListTitle = `Players [${playerList.length}/${maxPlayer}]`;
  console.log(maxPlayer);
  useEffect(() => {
    setPlayerListTitle(`Players [${playerList.length}/${maxPlayer}]`);
  }, [playerList.length, maxPlayer]);
  return (
    <div className="flex flex-col h-full ml-4 min-h-[70vh] max-h-[70vh]">
      <OnlineUserList
        name={playerListTitle}
        userData={userData}
        userList={playerList}
      />
      <HorizontalLine />
      <OnlineUserList
        name={"Watchers"}
        userData={userData}
        userList={watcherList}
      />
    </div>
  );
}
