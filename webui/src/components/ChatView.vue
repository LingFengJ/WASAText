<script>
export default {
    props: {
        id: {
            type: String,
            required: true
        }
    },
    data() {
        return {
            conversation: null,
            messages: [],
            newMessage: '',
            loading: true,
            error: null,
            authToken: null
        }
    },
    async created() {
        // this.authToken = sessionStorage.getItem('authToken');
        // await this.loadConversation();
        // this.scrollToBottom();

        console.log('Conversation ID:', this.id);
        this.authToken = sessionStorage.getItem('authToken');
        this.loadConversation();
    },
    methods: {
        async loadConversation() {
            try {
                const response = await fetch(`http://localhost:3000/conversations/${this.id}`, {
                    headers: {
                        'Authorization': `Bearer ${this.authToken}`
                    }
                });
                
                if (response.ok) {
                    const data = await response.json();
                    console.log('Conversation data:', data);
                    this.conversation = data.conversation;
                    this.messages = data.messages;
                } else {
                    this.error = 'Failed to load conversation';
                }
            } catch (error) {
                console.error('Error loading conversation:', error);
                this.error = 'Network error';
            } finally {
                this.loading = false;
            }
        },
        async sendMessage() {
            if (!this.newMessage.trim()) return;
            
            try {
                const response = await fetch(`http://localhost:3000/conversations/${this.id}/messages`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${this.authToken}`,
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        content: this.newMessage,
                        type: 'text'
                    })
                });
                
                if (response.ok) {
                    this.newMessage = '';
                    await this.loadConversation();
                    this.scrollToBottom();
                } else {
                    this.error = 'Failed to send message';
                }
            } catch (error) {
                console.error('Error sending message:', error);
                this.error = 'Network error';
            }
        },
        scrollToBottom() {
            this.$nextTick(() => {
                const container = this.$refs.messageContainer;
                if (container) {
                    container.scrollTop = container.scrollHeight;
                }
            });
        },
        getPhotoUrl(photoUrl) {
            if (!photoUrl) return null;
            const cleanPath = photoUrl.replace(/\\/g, '/');
            const normalizedPath = cleanPath.startsWith('/') ? cleanPath : '/' + cleanPath;
            return `http://localhost:3000${normalizedPath}`;
        }
    }
}
</script>

<template>
    <div class="chat-container">
        <!-- Conversation Header -->
        <div v-if="conversation" class="chat-header">
            <div class="d-flex align-items-center p-3 border-bottom">
                <div class="conversation-avatar me-3">
                    <img 
                        v-if="conversation.photoUrl" 
                        :src="getPhotoUrl(conversation.photoUrl)"
                        :alt="conversation.name"
                        class="avatar-img"
                    />
                    <i v-else class="bi" :class="conversation.type === 'group' ? 'bi-people-fill' : 'bi-person-circle'"></i>
                </div>
                <div>
                    <h5 class="mb-0">{{ conversation.name }}</h5>
                    <small>{{ conversation.type === 'group' ? 'Group Chat' : 'Direct Message' }}</small>
                </div>
            </div>
        </div>

        <!-- Loading State -->
        <div v-if="loading" class="text-center p-4">
            <div class="spinner-border" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
        
        <!-- Error State -->
        <div v-else-if="error" class="alert alert-danger m-3">
            {{ error }}
        </div>

        <!-- Messages -->
        <div v-else class="messages" ref="messageContainer">
            <div v-for="msg in messages" :key="msg.id" 
                 :class="['message', 'p-2', 'rounded', msg.senderId === authToken ? 'sent' : 'received']">
                <div class="message-content">
                    {{ msg.content }}
                </div>
                <div class="message-time small text-muted">
                    {{ new Date(msg.timestamp).toLocaleTimeString() }}
                </div>
            </div>
        </div>

        <!-- Message Input -->
        <div class="message-input">
            <form @submit.prevent="sendMessage">
                <div class="input-group">
                    <input 
                        type="text" 
                        class="form-control" 
                        v-model="newMessage" 
                        placeholder="Type a message..."
                    >
                    <button class="btn btn-primary" type="submit">Send</button>
                </div>
            </form>
        </div>
    </div>
</template>

<style scoped>
.chat-container {
    height: calc(100vh - 64px);
    display: flex;
    flex-direction: column;
    background-color: #f8f9fa;
}

.chat-header {
    background-color: white;
}

.messages {
    flex-grow: 1;
    overflow-y: auto;
    padding: 20px;
}

.message {
    margin: 10px 0;
    max-width: 70%;
}

.sent {
    margin-left: auto;
    background-color: #0d6efd;
    color: white;
}

.received {
    margin-right: auto;
    background-color: white;
}

.message-input {
    padding: 20px;
    background-color: white;
    border-top: 1px solid #ddd;
}

.conversation-avatar {
    width: 40px;
    height: 40px;
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