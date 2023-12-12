import { Room } from "@/types/playerTypes";
import { User } from "@/types/playerTypes";
function classNames(...classes: string[]): string {
  return classes.filter(Boolean).join(" ");
}
export default function RoomItemComponent({ room }: { room: Room }) {
  // console.log(room);
  return (
    <li key={room.room_id} className="border p-2 hover:bg-gray-100 rounded">
      <a href={`/room/${room.room_id}`}>
        <div className="mb-8 flex justify-between">
          <h3 className="text-2xl font-bold italic">{`Room: ${room.room_id}`}</h3>
          <h3
            className="text-lg font-semibold"
            style={{
              color: room.room_status === "open" ? "red" : "green",
            }}
          >
            {room.room_status}
          </h3>
        </div>
        <p
          className={`text-xl ${classNames(
            room.room_members && room.room_members.length > 0
              ? "text-gray-700"
              : "text-gray-500"
          )}`}
        >
          {room.room_members && room.room_members.length > 0
            ? room.room_members.map((member: User) => member.username) || []
            : "No members"}
        </p>

        {room.room_members &&
          room.room_members.length > 0 &&
          room.room_members.map((member: User) => (
            <div key={member.uuid}>{member.username}</div>
          ))}
        <p
          className="font-xl font-semibold"
          style={{
            color:
              room.game_status && room.game_status.isOver ? "red" : "green",
          }}
        >
          {room.game_status && room.game_status.isOver
            ? "Game Over"
            : "Game On"}
        </p>
      </a>
    </li>
  );
}
