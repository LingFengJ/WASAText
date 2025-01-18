<script>
export default {
    data() {
        return {
            conversations: [],
            loading: true,
            error: null
        }
    },
    async mounted() {
        await this.loadConversations();
    },
    methods: {
        getPhotoUrl(conv) {
            if (!conv.photoUrl) return null;
            // Fix Windows-style paths
            const cleanPath = conv.photoUrl.replace(/\\/g, '/');
            // Ensure the path starts with a single forward slash
            const normalizedPath = cleanPath.startsWith('/') ? cleanPath : '/' + cleanPath;
            return `http://localhost:3000${normalizedPath}`;
        },
        async loadConversations() {
            try {
                console.log('Loading conversations...');
                const response = await fetch('http://localhost:3000/conversations', {
                    headers: {
                        'Authorization': `Bearer ${sessionStorage.getItem('authToken')}`
                    }
                });
                
                console.log('Response status:', response.status);
                if (response.ok) {
                    const data = await response.json();
                    console.log('Conversations loaded:', data);
                    this.conversations = data;
                } else {
                    const errorData = await response.json();
                    this.error = errorData.message || 'Failed to load conversations';
                }
            } catch (error) {
                console.error('Error loading conversations:', error);
                this.error = 'Failed to load conversations';
            } finally {
                this.loading = false;
            }
        },
        formatTime(timestamp) {
            return new Date(timestamp).toLocaleTimeString();
        }
    }
}
</script>

<template>
    <div class="conversation-list">
        <div v-if="loading" class="text-center">
            <div class="spinner-border" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
        
        <div v-else-if="error" class="alert alert-danger">
            {{ error }}
        </div>

        <div v-else-if="conversations.length === 0" class="text-center p-4">
            <p>No conversations yet</p>
            <router-link to="/new-conversation" class="btn btn-primary">
                Start a conversation
            </router-link>
        </div>

        <div v-else class="list-group">
            <router-link 
                v-for="conv in conversations" 
                :key="conv.id"
                :to="`/conversations/${conv.id}`"
                class="list-group-item list-group-item-action"
            >
                           <!-- :src="`http://localhost:3000${conv.photoUrl}`" -->
                <div class="d-flex">
                    <div class="conversation-avatar me-3">
                        <img 
                            v-if="conv.photoUrl"  
                            
                            :src="getPhotoUrl(conv)"
                            :alt="conv.name"
                            class="avatar-img"
                        >
                        <i v-else class="bi" :class="conv.type === 'group' ? 'bi-people-fill' : 'bi-person-circle'"></i>
                    </div>
                    <div class="flex-grow-1">
                        <div class="d-flex w-100 justify-content-between">
                            <h5 class="mb-1">{{ conv.name }}</h5>
                            <small v-if="conv.lastMessage">
                                {{ formatTime(conv.lastMessage.timestamp) }}
                            </small>
                        </div>
                        <p class="mb-1" v-if="conv.lastMessage">
                            <i v-if="conv.lastMessage.type === 'photo'" class="bi bi-image me-1"></i>
                            {{ conv.lastMessage.content }}
                        </p>
                        <small>{{ conv.type === 'group' ? 'Group' : 'Direct Message' }}</small>
                    </div>
                </div>
            </router-link>
        </div>
    </div>
</template>

<style scoped>
.conversation-list {
    padding: 20px;
}

.list-group-item {
    margin-bottom: 8px;
    border-radius: 8px;
}

.list-group-item:hover {
    background-color: #f8f9fa;
}

.conversation-avatar {
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.avatar-img {
    width: 100%;
    height: 100%;
    border-radius: 50%;
    object-fit: cover;
}

.bi-person-circle,
.bi-people-fill {
    font-size: 2rem;
    color: #6c757d;
}
</style>


