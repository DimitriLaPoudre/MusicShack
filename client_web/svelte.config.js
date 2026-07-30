import adapter from '@sveltejs/adapter-static';
import {preprocess} from 'svelte/compiler';
export default {
  kit: {
    adapter: adapter({
      pages: 'build',
      assets: 'build',
      fallback: 'index.html'
    })
  },
};
