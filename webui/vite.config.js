import {fileURLToPath, URL} from 'node:url'
import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig(({command, mode}) => {
    const ret = {
        plugins: [vue()],
        resolve: {
            alias: {
                '@': fileURLToPath(new URL('./src', import.meta.url))
            }
        },
    };
    
    // Set API URL based on environment
    let apiUrl = "http://localhost:3000"; // Default for development
    
    if (mode === 'production') {
        // Use the Railway backend URL in production
        apiUrl = "https://wasatext-production.up.railway.app"; // Replace with your actual Railway URL
    }
    
    ret.define = {
        "__API_URL__": JSON.stringify(apiUrl),
    };
    
    return ret;
})