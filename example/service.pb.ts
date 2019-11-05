
export interface Hat {
  size: number;
  color: string;
  name: string;
  ribons?: Ribon[];
}

export interface HatJSON {
  size?: number;
  color?: string;
  name?: string;
  ribons?: RibonJSON[];
}

export const HatToJSON = (m: Hat): HatJSON => {
  return {
    size: m.size,
    color: m.color,
    name: m.name,
    ribons: m.ribons && m.ribons.map(RibonToJSON),
  };
};

export const JSONToHat = (m: HatJSON): Hat => {
  return {
    size: m.size || 0,
    color: m.color || "",
    name: m.name || "",
    ribons: m.ribons && m.ribons.map(JSONToRibon),
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
