"use client";

import About from "@/Components/About";
import Footer from "@/Components/Footer";
import { mapColor, Obstacles } from "@/utils/MapInfo";
import { SimulationStateSchema } from "@/utils/Schemas";
import {
  detectPathDirectionChange,
  driverNameMapping,
  formatNumber,
  PathDirection,
} from "@/utils/Utils";
import React, { useEffect, useState } from "react";
import { z } from "zod";

const SQUARE_SIZE = 20;

export default function Home() {
  const [currentSimulationState, setCurrentSimulationState] = useState<
    z.infer<typeof SimulationStateSchema>
  >([]);

  useEffect(() => {
    const socket = new WebSocket(
      process.env.NEXT_PUBLIC_BACKEND_WS_URL as string
    );

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);

      const validatedData = SimulationStateSchema.safeParse(data);
      if (validatedData.success) {
        setCurrentSimulationState(validatedData.data);
      } else {
        console.log(validatedData.error);
      }
    };
  }, []);

  const mapRenderHelper = () => {
    const nodes: React.JSX.Element[] = [];

    // Builds empty map
    for (let y = 0; y < 50; y++) {
      for (let x = 0; x < 50; x++) {
        nodes.push(
          <rect
            key={`${x}-${y}`}
            height={SQUARE_SIZE}
            width={SQUARE_SIZE}
            fill={mapColor["road"]}
            x={x * SQUARE_SIZE}
            y={y * SQUARE_SIZE}
          />
        );
      }
    }

    // Places obstacles such as park, river and buildings
    Obstacles.forEach((obstacle) => {
      const { xStart, xEnd, yStart, yEnd, type } = obstacle;

      nodes.push(
        <rect
          key={`${xStart}-${yStart}-obstacle`}
          fill={mapColor[type]}
          x={xStart * SQUARE_SIZE}
          y={yStart * SQUARE_SIZE}
          width={Math.abs(xEnd - xStart) * SQUARE_SIZE}
          height={Math.abs(yEnd - yStart) * SQUARE_SIZE}
        />
      );
    });

    // Places the position of drivers, path to the pickUp/dropOff point and the markers for pickUp/dropOff point
    currentSimulationState.forEach((driver) => {
      if (driver.path.length === 0) return;
      let pathDirection: PathDirection = detectPathDirectionChange(
        driver.currentCoordinates,
        driver.path[0]
      );
      for (let i = 0; i < driver.path.length; i++) {
        const CoordInfo = driver.path[i];
        let nextPathDirection = pathDirection;
        if (driver.path.length > 1 && i != driver.path.length - 1) {
          const nextCoordInfo = driver.path[i + 1];
          nextPathDirection = detectPathDirectionChange(
            CoordInfo,
            nextCoordInfo
          );
        }

        if (nextPathDirection != pathDirection) {
          const calculatedSquareSize = SQUARE_SIZE - SQUARE_SIZE / 4;

          nodes.push(
            <rect
              key={`${driver.driver}-path-${CoordInfo.X}-${CoordInfo.Y}-firstHalf`}
              x={
                pathDirection === "Right"
                  ? CoordInfo.X * SQUARE_SIZE
                  : CoordInfo.X * SQUARE_SIZE + SQUARE_SIZE / 4
              }
              y={
                pathDirection === "Down"
                  ? CoordInfo.Y * SQUARE_SIZE
                  : CoordInfo.Y * SQUARE_SIZE + SQUARE_SIZE / 4
              }
              width={
                pathDirection == "Right"
                  ? calculatedSquareSize
                  : pathDirection == "Left"
                  ? calculatedSquareSize
                  : SQUARE_SIZE / 2
              }
              height={
                pathDirection == "Up"
                  ? calculatedSquareSize
                  : pathDirection == "Down"
                  ? calculatedSquareSize
                  : SQUARE_SIZE / 2
              }
              fill={mapColor["path"]}
            />,
            <rect
              key={`${driver.driver}-path-${CoordInfo.X}-${CoordInfo.Y}-secondHalf`}
              x={
                nextPathDirection === "Left"
                  ? CoordInfo.X * SQUARE_SIZE
                  : CoordInfo.X * SQUARE_SIZE + SQUARE_SIZE / 4
              }
              y={
                nextPathDirection === "Up"
                  ? CoordInfo.Y * SQUARE_SIZE
                  : CoordInfo.Y * SQUARE_SIZE + SQUARE_SIZE / 4
              }
              width={
                nextPathDirection == "Right"
                  ? calculatedSquareSize
                  : nextPathDirection == "Left"
                  ? calculatedSquareSize
                  : SQUARE_SIZE / 2
              }
              height={
                nextPathDirection == "Up"
                  ? calculatedSquareSize
                  : nextPathDirection == "Down"
                  ? calculatedSquareSize
                  : SQUARE_SIZE / 2
              }
              fill={mapColor["path"]}
            />
          );

          pathDirection = nextPathDirection;
        } else {
          nodes.push(
            <rect
              key={`${driver.driver}-path-${CoordInfo.X}-${CoordInfo.Y}`}
              x={
                pathDirection === "Right"
                  ? CoordInfo.X * SQUARE_SIZE
                  : pathDirection === "Left"
                  ? CoordInfo.X * SQUARE_SIZE
                  : CoordInfo.X * SQUARE_SIZE + SQUARE_SIZE / 4
              }
              y={
                pathDirection === "Up"
                  ? CoordInfo.Y * SQUARE_SIZE
                  : pathDirection === "Down"
                  ? CoordInfo.Y * SQUARE_SIZE
                  : CoordInfo.Y * SQUARE_SIZE + SQUARE_SIZE / 4
              }
              width={
                pathDirection == "Right"
                  ? SQUARE_SIZE
                  : pathDirection == "Left"
                  ? SQUARE_SIZE
                  : SQUARE_SIZE / 2
              }
              height={
                pathDirection == "Up"
                  ? SQUARE_SIZE
                  : pathDirection == "Down"
                  ? SQUARE_SIZE
                  : SQUARE_SIZE / 2
              }
              fill={mapColor["path"]}
            />
          );
        }
      }
    });

    currentSimulationState.forEach((driver) => {
      if (driver.driverStatus == "Picking Up") {
        nodes.push(
          <circle
            key={`${driver.driver}-pickUpCoordinates`}
            fill={mapColor["pickUpPoint"]}
            cx={driver.pickUpCoordinates.X * SQUARE_SIZE + SQUARE_SIZE / 2}
            cy={driver.pickUpCoordinates.Y * SQUARE_SIZE + SQUARE_SIZE / 2}
            r={SQUARE_SIZE / 2}
          />
        );
      } else if (driver.driverStatus === "Dropping Off") {
        nodes.push(
          <circle
            key={`${driver.driver}-dropOffCoordinates`}
            fill={mapColor["dropOffPoint"]}
            cx={driver.dropOffCoordinates.X * SQUARE_SIZE + SQUARE_SIZE / 2}
            cy={driver.dropOffCoordinates.Y * SQUARE_SIZE + SQUARE_SIZE / 2}
            r={SQUARE_SIZE / 2}
          />
        );
      }

      nodes.push(
        <circle
          key={`${driver.driver}-currentCoordinates`}
          fill={mapColor["driverLocation"]}
          cx={driver.currentCoordinates.X * SQUARE_SIZE + SQUARE_SIZE / 2}
          cy={driver.currentCoordinates.Y * SQUARE_SIZE + SQUARE_SIZE / 2}
          r={SQUARE_SIZE / 2}
        />
      );
    });

    return (
      <svg viewBox="0 0 1000 1000" className="h-full w-full">
        {nodes}
      </svg>
    );
  };

  const driverListRenderHelper = currentSimulationState.map((driver) => {
    return (
      <div
        className="border-[1px] border-border rounded p-4 text-text-primary bg-bg-secondary flex flex-col gap-1"
        key={`${driver.driver}-details`}
      >
        <div className="flex justify-between items-center">
          <p className="text-text-secondary/80 uppercase font-semibold text-sm">
            Name
          </p>
          <p>{driverNameMapping[driver.driver]}</p>
        </div>
        <div className="flex justify-between items-center">
          <p className="text-text-secondary/80 uppercase font-semibold text-sm">
            Status
          </p>
          <p>{driver.driverStatus}</p>
        </div>
        <div className="flex justify-between items-center">
          <p className="text-text-secondary/80 uppercase font-semibold text-sm">
            PickUp point
          </p>
          <p className="tabular-nums">{`${formatNumber(
            driver.pickUpCoordinates.X
          )}, ${formatNumber(driver.pickUpCoordinates.Y)}`}</p>
        </div>
        <div className="flex justify-between items-center">
          <p className="text-text-secondary/80 uppercase font-semibold text-sm">
            DropOff point
          </p>
          <p className="tabular-nums">{`${formatNumber(
            driver.dropOffCoordinates.X
          )}, ${formatNumber(driver.dropOffCoordinates.Y)}`}</p>
        </div>
      </div>
    );
  });

  return (
    <div className=" md:flex min-h-screen md:h-screen justify-center items-center bg-bg-primary text-text-primary">
      <div className="h-full">{mapRenderHelper()}</div>
      <div className="grow md:min-w-[400px] h-full flex flex-col justify-between border-l-[1px] border-l-border">
        <div className="px-8 grow flex flex-col shrink">
          <About />
          <div className="grow flex flex-col mb-8">
            <h2 className="text-lg mb-2 font-semibold">Drivers</h2>
            <div className="grow flex flex-col gap-2 overflow-auto h-full md:h-0">
              {driverListRenderHelper}
            </div>
          </div>
        </div>
        <Footer />
      </div>
    </div>
  );
}
