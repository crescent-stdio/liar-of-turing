import * as datetime from "date-fns";

export const getDayFromTimestamp = (timestamp: number): string => {
  return datetime.format(new Date(timestamp), "yyyy-MM-dd HH:mm");
};
export const getTimeFromTimer = (time: number): string => {
  //using date-fns ans 10:00
  const date = new Date(time * 1000);
  let minutes = (date.getMinutes() < 10 ? "0" : "") + date.getMinutes();
  let seconds = (date.getSeconds() < 10 ? "0" : "") + date.getSeconds();
  return `${minutes}:${seconds}`;
};
