import { Feather, FeatherJSON, FeatherToJSON, JSONToFeather } from '../drawer/feather.pb'
// @@protoc_insertion_point(plugin_imports)

export interface Hat {
  size: number;
  color: string;
  name: string;
  ribons?: Ribon[];
  plume?: HatPlumeEntry;
  createTime?: string;
}

export interface HatJSON {
  size?: number;
  color?: string;
  name?: string;
  ribons?: RibonJSON[];
  plume?: HatPlumeEntryJSON;
  create_time?: string;
}

export const HatToJSON = (m: Hat): HatJSON => {
  return {
    size: m.size,
    color: m.color,
    name: m.name,
    ribons: m.ribons && m.ribons.map(RibonToJSON),
    plume: m.plume && HatPlumeEntryToJSON(m.plume),
    create_time: m.createTime,
  };
};

export const JSONToHat = (m: HatJSON): Hat => {
  return {
    size: m.size || 0,
    color: m.color || "",
    name: m.name || "",
    ribons: m.ribons && m.ribons.map(JSONToRibon),
    plume: m.plume && JSONToHatPlumeEntry(m.plume),
    createTime: m.create_time,
  };
};
export interface Size {
  inches: number;
}

export interface SizeJSON {
  inches?: number;
}

export const SizeToJSON = (m: Size): SizeJSON => {
  return {
    inches: m.inches,
  };
};

export const JSONToSize = (m: SizeJSON): Size => {
  return {
    inches: m.inches || 0,
  };
};
export interface Ribon {
  color: string;
}

export interface RibonJSON {
  color?: string;
}

export const RibonToJSON = (m: Ribon): RibonJSON => {
  return {
    color: m.color,
  };
};

export const JSONToRibon = (m: RibonJSON): Ribon => {
  return {
    color: m.color || "",
  };
};
export interface HatPlumeEntry {
  [key: string]: Feather;
}

export interface HatPlumeEntryJSON {
  [key: string]: FeatherJSON;
}
export const JSONToHatPlumeEntry = (m: HatPlumeEntryJSON): HatPlumeEntry => {
  return Object.keys(m).reduce((acc, key) => {
    acc[key] = JSONToFeather(m[key]);
    return acc;
  }, {} as HatPlumeEntry);
};

export const HatPlumeEntryToJSON = (m: HatPlumeEntry): HatPlumeEntryJSON => {
  return Object.keys(m).reduce((acc, key) => {
    acc[key] = FeatherToJSON(m[key]);
    return acc;
  }, {} as HatPlumeEntryJSON);
};
