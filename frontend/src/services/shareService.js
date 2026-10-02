import { fileService } from './fileService.js';

export const shareService = {
  list: () => fileService.listShares(),
  create: (...args) => fileService.createShare(...args),
  revoke: (id) => fileService.revokeShare(id),
  get: async (id) => {
    const shares = await fileService.listShares();
    return shares.find((s) => (s.id === id || s.token === id || s.shareId === id) && s.active) || null;
  },
};
