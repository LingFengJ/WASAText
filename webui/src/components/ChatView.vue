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
            authToken: null,
            username: null,
            showingReactionModal: false,
            showingForwardModal: false,
            selectedMessage: null,
            conversations: [],
            emojis: ['👍', '❤️', '😊', '😂', '😮', '😢']
        }
    },
    async created() {
        console.log('Conversation ID:', this.id);
        this.authToken = sessionStorage.getItem('authToken');
        this.username = sessionStorage.getItem('username');
        console.log('Current username:', this.username);
        this.loadConversation();
    },
    methods: {
        async loadConversation() {
            try {
                const response = await this.$axios.get(`/conversations/${this.id}`, {
                    headers: {
                        'Authorization': `Bearer ${this.authToken}`
                    }
                });
                
                const data = response.data;
                console.log('Conversation data:', data);
                this.conversation = data.conversation;
                // this.messages = data.messages;
                this.messages = data.messages.map(msg => ({
                    ...msg,
                    reactions: msg.reactions || [],
                    // Status should update based on backend response
                    status: msg.status || 'sent'
                }));
                console.log('messages with reactions:', this.messages);

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
                const response = await this.$axios.post(`/conversations/${this.id}/messages`, 
                    {
                    content: this.newMessage,
                    type: 'text'
                    },    
                    {                
                    headers: {
                        'Authorization': `Bearer ${this.authToken}`,
                        'Content-Type': 'application/json'
                    },
                    }
                );
                
                this.newMessage = '';
                await this.loadConversation();
                this.scrollToBottom();
            } catch (error) {
                if (error.response) {
                    const statusCode = error.response.status;
                    switch (statusCode) {
                        case 400:
                            console.error('Bad request');
                        case 401:
                            console.error('Access Unauthorized:', error.response.data);
                            break;
                        case 403:
                            console.error('Access Forbidden: You are not a member of the conversation', error.response.data);
                            break;
                        case 404:
                            console.error('Recipient Not Found:', error.response.data);
                            break;
                        case 500:
                            console.error('Failed to get Coversation Internal Server Error:', error.response.data);
                            break;
                        default:
                            console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
                    }
                } else {
                        console.error('Error sending message:', error);
                        this.error = 'Network error';
                    }
            
            }
        },

        showReactionPicker(message) {
            if (this.selectedMessage?.id === message.id) {
                this.selectedMessage = null;
                this.showingReactionModal = false;
            } else {
                this.selectedMessage = message;
                this.showingReactionModal = true;
            }
        },

        // Click outside handler to close reaction picker
        closeReactionPicker(event) {
            if (!event.target.closest('.reaction-modal') && !event.target.closest('.message-actions')) {
                this.showingReactionModal = false;
                this.selectedMessage = null;
            }
        }, 
        async addReaction(message, emoji) {
            try {
                const response = await this.$axios.post(`/messages/${message.id}/reactions`, 
                {
                    Emoji:  emoji
                },
                {
                    headers: {
                        'Authorization': `Bearer ${this.authToken}`,
                        'Content-Type': 'application/json'
                    }
                });
                
                await this.loadConversation();
            } catch (error) {
                if (err.response){    
                    const statusCode = error.response.status;
                    switch (statusCode) {
                        case 400:
                            console.error('Bad request');
                        case 401:
                            console.error('Access Unauthorized:', error.response.data);
                            break;
                        case 403:
                            console.error('Access Forbidden: You are not a member of the conversation', error.response.data);
                            break;
                        case 404:
                            console.error('Recipient Not Found:', error.response.data);
                            break;
                        case 409:
                            console.error('Conflict:', error.response.data);
                            break;
                        case 500:
                            console.error('Failed to get Coversation Internal Server Error:', error.response.data);
                            break;
                        default:
                            console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
                    }
                }else{
                    console.error('Error adding reaction:', error);
                }
            }
            this.showingReactionModal = false;
            this.selectedMessage = null;
        },

        async removeReaction(messageId, reaction) {
            try {
                const response = await this.$axios.delete(`/messages/${messageId}/reactions`, {
                    headers: {
                        'Authorization': `Bearer ${this.authToken}`
                    }
                });
       await this.loadConversation();
            } catch (error) {
                console.error('Error removing reaction:', error);
            }
        },
        showForwardDialog(message) {
            console.log('Opening forward dialog for message:', message);
            this.selectedMessage = message;
            this.showingForwardModal = true;
            this.loadConversations();
        },
        // openForwardDialog(message) {
        //     this.selectedMessage = message;
        //     this.showForwardModal = true;
        //     this.loadConversations();
        //     },
        closeForwardDialog() {
            this.showingForwardModal = false;
            this.selectedMessage = null;
            },
        async loadConversations() {
            try {
                const response = await this.$axios.get('/conversations', {
                    headers: {
                        'Authorization': `Bearer ${this.authToken}`
                    }
                });
                
                this.conversations = response.data;
                
            } catch (error) {
                console.error('Error loading conversations:', error);
            }
        },
        async forwardMessage(message, targetconversationId) {
            try {
                const response = await this.$axios.post(`/messages/${message.id}/forward`, 
                    {
                        conversationId: targetconversationId
                    },
                    {
                        headers: {
                            'Authorization': `Bearer ${this.authToken}`
                        }
                    }
                );
                
                // this.showingForwardModal = false;
                this.closeForwardDialog();
            } catch (error) {
                console.error('Error forwarding message:', error);
            }
        },

        async deleteMessage(messageId) {
            try {
                await this.$axios.delete(`/messages/${messageId}`, {
                headers: {
                    'Authorization': `Bearer ${this.authToken}`
                }
                });
                await this.loadConversation();
            } catch (error) {
                console.error('Error deleting message:', error);
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
            // return `http://localhost:3000${normalizedPath}`;
            return `${__API_URL__}${normalizedPath}`;
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
    <div class="messages" ref="messageContainer">
        <div v-for="msg in messages" :key="msg.id" 
            class = "message p-2 rounded mb-2"
            :class="msg.senderId === authToken ? 'sent' : 'received'">
        <div class="message-content">
            {{ msg.content }}

        <!-- Message status for sent messages -->
        <span v-if="msg.senderUsername === username" class="message-status">
            <i class="bi" :class="msg.status === 'read' ? 'bi-check2-all' : 'bi-check2'"></i>
        </span>            
        <!-- Sender username for received messages -->
        <div v-else class="message-sender">
            {{ msg.senderUsername }}
        </div>
        </div>



    <!-- Show existing reactions -->
    <div v-if="msg.reactions && msg.reactions.length > 0" 
         class="message-reactions mt-2">
      <span v-for="reaction in msg.reactions" 
            :key="reaction.userId + reaction.emoji"
            class="reaction"
            @click="removeReaction(msg.id, reaction)">
        {{ reaction.emoji }}
        <span class="reaction-username">{{ reaction.username }}</span>
      </span>
    </div>

    <div v-if="showingReactionModal && selectedMessage?.id === msg.id" 
        class="reaction-picker-container">
        <div class="reaction-picker">
        <span v-for="emoji in emojis" 
             :key="emoji"
             class="emoji"
             @click.stop="addReaction(msg, emoji)">
          {{ emoji }}
        </span>
      </div>
    </div>

    <div class="message-footer d-flex justify-content-between align-items-center mt-1">
        <div class="message-time small text-muted">
            {{ new Date(msg.timestamp).toLocaleTimeString() }}
        </div>
        <!-- Message actions (emoji and forward buttons) -->
        <div class="message-actions">
            <button class="btn btn-sm btn-link" @click.stop="showReactionPicker(msg)">
                <i class="bi bi-emoji-smile"></i>
            </button>
            <button class="btn btn-sm btn-link" @click.stop="showForwardDialog(msg)">
                <i class="bi bi-forward"></i>
            </button>

            <button class="btn btn-sm btn-link" @click.stop="deleteMessage(msg.id)">
                <i class="bi bi-trash"></i>
            </button>
        </div>
    </div>
    </div>
    </div>

    <!-- Message input -->
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

    <!-- Reaction Picker Modal -->
    <div v-if="showingReactionModal" class="reaction-modal">
        <div class="reaction-picker">
            <div v-for="emoji in emojis" :key="emoji" 
                    class="emoji" @click="addReaction(selectedMessage, emoji)">
                {{ emoji }}
            </div>
        </div>
    </div>

  <!-- Forward Dialog -->
  <div v-if="showingForwardModal" class="modal" style="display: block">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Forward Message</h5>
          <button type="button" class="btn-close" @click="closeForwardDialog"></button>
        </div>
        <div class="modal-body">
          <div class="list-group">
            <button v-for="conv in conversations" 
                    :key="conv.id"
                    class="list-group-item list-group-item-action"
                    @click="forwardMessage(selectedMessage, conv.id)">
              {{ conv.name }}
            </button>
          </div>
        </div>
      </div>
    </div>
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

.messages {
    flex-grow: 1;
    overflow-y: auto;
    padding: 20px;
}

.message {
    position: relative;
    margin: 10px 0;
    max-width: 70%;
}

.message-input {
    position: sticky;
    bottom: 0;
    padding: 20px;
    background-color: white;
    border-top: 1px solid #ddd;
    margin-top: auto;
}

.message-content {
    word-break: break-word;
    white-space: pre-wrap;
    max-width: 200%;
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

.message-actions {
    position: absolute;
    right: -60px;
    top: 50%;
    transform: translateY(-50%);
    display: none;
    background: white;
    border-radius: 20px;
    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
    padding: 4px;
}

.message:hover .message-actions {
    display: flex;
}

.reaction-modal {
    position: absolute;
    bottom: 100%;
    left: 0;
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    z-index: 1000;
    padding: 8px;
    margin-bottom: 8px;
}

.reaction-picker {
    display: flex;
    gap: 8px;
}

.emoji {
    font-size: 1.2em;
    cursor: pointer;
    padding: 4px;
    border-radius: 4px;
    transition: background-color 0.2s;
}

.emoji:hover {
    background: #f0f0f0;
}

.message-reactions {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    margin-top: 4px;
}

.reaction {
    display: inline-flex;
    align-items: center;
    background: rgba(0, 0, 0, 0.05);
    border-radius: 12px;
    padding: 2px 8px;
    font-size: 0.9em;
    cursor: pointer;
}

.reaction-username {
    font-size: 0.8em;
    margin-left: 4px;
    color: #666;
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

.message-status {
    display: inline-block;   
    color: #8e8e8e;
    position: absolute;
    bottom: 4px;
    right: 8px;
}

.message-sender {
    font-size: 0.8em;
    color: #666;
    margin-bottom: 2px;
}

.sent .message-status {
    color: #fff;
}
</style>