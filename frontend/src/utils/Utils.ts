export type PathDirection = "Up" | "Down" | "Left" | "Right";

export const detectPathDirectionChange = (
  curr: { X: number; Y: number },
  next: { X: number; Y: number }
): PathDirection => {
  if (curr.X < next.X) return "Right";
  else if (curr.X > next.X) return "Left";
  else if (curr.Y < next.Y) return "Down";
  else return "Up";
};

export const driverNameMapping: Record<string, string> = {
  A: "Amy",
  B: "Bretha",
  C: "Clair",
  D: "Daina",
  E: "Emily",
};

export const formatNumber = (num: number) => {
  return num < 10 ? "0" + num : num;
};
