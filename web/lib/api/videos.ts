import type {
  VideoDownloadURLRequest,
  VideoDownloadURLResponse,
  VideoUploadURLRequest,
  VideoUploadURLResponse,
} from '@/types/api';
import { api } from './client';

export const videosApi = {
  async getUploadUrl(data: VideoUploadURLRequest): Promise<VideoUploadURLResponse> {
    return api.post('videos/upload-url', { json: data }).json<VideoUploadURLResponse>();
  },

  async getDownloadUrl(data: VideoDownloadURLRequest): Promise<VideoDownloadURLResponse> {
    return api.post('videos/download-url', { json: data }).json<VideoDownloadURLResponse>();
  },

  /** Upload a file directly to the pre-signed URL (bypasses ky, no auth header) */
  async uploadToPresignedUrl(uploadUrl: string, file: File): Promise<void> {
    const response = await fetch(uploadUrl, {
      method: 'PUT',
      headers: { 'Content-Type': file.type },
      body: file,
    });
    if (!response.ok) {
      throw new Error(`Upload failed: ${response.status} ${response.statusText}`);
    }
  },
};
