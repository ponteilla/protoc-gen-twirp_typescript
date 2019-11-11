// @@protoc_insertion_point(plugin_imports)

export interface Glitter {
  isGel: boolean;
}

export interface GlitterJSON {
  is_gel?: boolean;
}

export const GlitterToJSON = (m: Glitter): GlitterJSON => {
  return {
    is_gel: m.isGel,
  };
};

export const JSONToGlitter = (m: GlitterJSON): Glitter => {
  return {
    isGel: m.is_gel || false,
  };
};
