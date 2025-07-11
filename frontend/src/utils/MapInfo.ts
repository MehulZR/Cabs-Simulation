type Obstacle = {
  xStart: number;
  xEnd: number;
  yStart: number;
  yEnd: number;
  type: "park" | "building" | "river";
};

export const Obstacles: Obstacle[] = [
  { xStart: 0, xEnd: 5, yStart: 1, yEnd: 7, type: "building" },
  { xStart: 40, xEnd: 49, yStart: 17, yEnd: 23, type: "park" },
  { xStart: 6, xEnd: 9, yStart: 0, yEnd: 5, type: "building" },
  { xStart: 6, xEnd: 9, yStart: 6, yEnd: 10, type: "building" },
  { xStart: 10, xEnd: 16, yStart: 0, yEnd: 10, type: "building" },
  { xStart: 17, xEnd: 18, yStart: 0, yEnd: 2, type: "building" },
  { xStart: 17, xEnd: 18, yStart: 3, yEnd: 7, type: "building" },
  { xStart: 17, xEnd: 18, yStart: 8, yEnd: 12, type: "building" },
  { xStart: 19, xEnd: 25, yStart: 1, yEnd: 12, type: "park" },
  { xStart: 26, xEnd: 29, yStart: 0, yEnd: 13, type: "river" },
  { xStart: 26, xEnd: 29, yStart: 14, yEnd: 15, type: "river" },
  { xStart: 26, xEnd: 29, yStart: 16, yEnd: 28, type: "river" },
  { xStart: 0, xEnd: 15, yStart: 28, yEnd: 30, type: "river" },
  { xStart: 16, xEnd: 17, yStart: 28, yEnd: 30, type: "river" },
  { xStart: 18, xEnd: 40, yStart: 28, yEnd: 30, type: "river" },
  { xStart: 39, xEnd: 41, yStart: 28, yEnd: 50, type: "river" },
  { xStart: 0, xEnd: 5, yStart: 8, yEnd: 10, type: "park" },
  { xStart: 0, xEnd: 16, yStart: 11, yEnd: 15, type: "building" },
  { xStart: 16, xEnd: 17, yStart: 13, yEnd: 15, type: "building" },
  { xStart: 18, xEnd: 25, yStart: 13, yEnd: 15, type: "building" },
  { xStart: 0, xEnd: 14, yStart: 16, yEnd: 27, type: "park" },
  { xStart: 15, xEnd: 22, yStart: 16, yEnd: 27, type: "building" },
  { xStart: 23, xEnd: 25, yStart: 16, yEnd: 27, type: "building" },
  { xStart: 30, xEnd: 50, yStart: 0, yEnd: 2, type: "building" },
  { xStart: 30, xEnd: 39, yStart: 3, yEnd: 6, type: "building" },
  { xStart: 40, xEnd: 50, yStart: 3, yEnd: 6, type: "building" },
  { xStart: 30, xEnd: 34, yStart: 7, yEnd: 10, type: "building" },
  { xStart: 35, xEnd: 43, yStart: 7, yEnd: 10, type: "park" },
  { xStart: 44, xEnd: 50, yStart: 7, yEnd: 10, type: "building" },
  { xStart: 30, xEnd: 48, yStart: 11, yEnd: 16, type: "building" },
  { xStart: 49, xEnd: 50, yStart: 11, yEnd: 16, type: "building" },
  { xStart: 30, xEnd: 34, yStart: 17, yEnd: 27, type: "building" },
  { xStart: 35, xEnd: 39, yStart: 17, yEnd: 24, type: "building" },
  { xStart: 35, xEnd: 39, yStart: 25, yEnd: 27, type: "park" },
  { xStart: 40, xEnd: 50, yStart: 24, yEnd: 27, type: "building" },
  { xStart: 42, xEnd: 49, yStart: 28, yEnd: 31, type: "building" },
  { xStart: 42, xEnd: 49, yStart: 32, yEnd: 42, type: "building" },
  { xStart: 42, xEnd: 49, yStart: 43, yEnd: 50, type: "park" },
  { xStart: 1, xEnd: 15, yStart: 31, yEnd: 40, type: "building" },
  { xStart: 1, xEnd: 15, yStart: 41, yEnd: 44, type: "park" },
  { xStart: 0, xEnd: 15, yStart: 45, yEnd: 50, type: "building" },
  { xStart: 29, xEnd: 38, yStart: 31, yEnd: 41, type: "park" },
  { xStart: 29, xEnd: 38, yStart: 42, yEnd: 44, type: "building" },
  { xStart: 29, xEnd: 33, yStart: 45, yEnd: 49, type: "building" },
  { xStart: 34, xEnd: 38, yStart: 45, yEnd: 49, type: "building" },
  { xStart: 23, xEnd: 28, yStart: 31, yEnd: 41, type: "building" },
  { xStart: 24, xEnd: 28, yStart: 42, yEnd: 46, type: "building" },
  { xStart: 23, xEnd: 28, yStart: 47, yEnd: 50, type: "building" },
  { xStart: 16, xEnd: 23, yStart: 42, yEnd: 50, type: "building" },
  { xStart: 16, xEnd: 22, yStart: 31, yEnd: 34, type: "building" },
  { xStart: 16, xEnd: 22, yStart: 35, yEnd: 41, type: "park" },
];

export const mapColor = {
  park: "#065f46",
  river: "#1e40af",
  building: "#0f172a",
  pickUpPoint: "#be123c",
  dropOffPoint: "#16a34a",
  driverLocation: "#0f172a",
  path: "#64748b",
  road: "#334155",
};
