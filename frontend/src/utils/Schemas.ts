import { z } from "zod";

export const CoordinateSchema = z.object({
  X: z.number(),
  Y: z.number(),
});

export const SimulationStateSchema = z
  .object({
    driver: z.string(),
    driverStatus: z.enum(["Picking Up", "Dropping Off", "Available"]),
    currentCoordinates: CoordinateSchema,
    pickUpCoordinates: CoordinateSchema,
    dropOffCoordinates: CoordinateSchema,
    path: CoordinateSchema.array(),
  })
  .array();
