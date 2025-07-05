import { mapColor } from "@/utils/MapInfo";

const About = () => {
  return (
    <>
      <div className="my-8 shrink">
        <h2 className="text-text-primary text-lg font-semibold">
          Cabs Simulation
        </h2>
        <p className="text-text-secondary">
          An event driven simulation of Cabs/Delivery partners in a 50 x 50 map.
        </p>
      </div>
      <div className="mb-8">
        <h2 className="text-lg mb-2 font-semibold">Map Legend</h2>
        <ul className="flex flex-col px-4 gap-1">
          <li className="flex gap-4 items-center">
            <svg width={16} height={16}>
              <circle r={8} cx={8} cy={8} fill={mapColor["pickUpPoint"]} />
            </svg>
            PickUp Point
          </li>
          <li className="flex gap-4 items-center">
            <svg width={16} height={16}>
              <circle r={8} cx={8} cy={8} fill={mapColor["dropOffPoint"]} />
            </svg>
            DropOff Point
          </li>
          <li className="flex gap-4 items-center">
            <svg width={16} height={16}>
              <circle r={8} cx={8} cy={8} fill={mapColor["driverLocation"]} />
            </svg>
            Driver&apos;s Location
          </li>
        </ul>
      </div>
    </>
  );
};

export default About;
