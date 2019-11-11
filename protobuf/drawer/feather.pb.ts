

import { Glitter, GlitterJSON, GlitterToJSON, JSONToGlitter } from './glitter.pb'

// @@protoc_insertion_point(plugin_imports)

export interface Feather {
  size: number;
  color: string;
  glitter?: Glitter[];
}

export interface FeatherJSON {
  size?: number;
  color?: string;
  glitter?: GlitterJSON[];
}

export const FeatherToJSON = (m: Feather): FeatherJSON => {
  return {
    size: m.size,
    color: m.color,
    glitter: m.glitter && m.glitter.map(GlitterToJSON),
  };
};

export const JSONToFeather = (m: FeatherJSON): Feather => {
  return {
    size: m.size || 0,
    color: m.color || "",
    glitter: m.glitter && m.glitter.map(JSONToGlitter),
  };
};