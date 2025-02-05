<script>
export default {
    data() {
        return {
            recipientName: '',
            message: '',
            error: null,
            loading: false
        }
    },
    methods: {
        async startConversation() {
            if (!this.recipientName || !this.message) {
                this.error = 'Please fill in all fields';
                return;
            }

            this.loading = true;
            try {
                const response = await this.$axios.post('/conversations//messages',                     
                {
                    recipientName: this.recipientName,
                    content: this.message,
                    type: 'text'
                }, {
                    headers: {
                        'Authorization': `Bearer ${sessionStorage.getItem('authToken')}`,
                        'Content-Type': 'application/json'
                    },

                });

                const data = response.data;
                this.$router.push(`/conversations/${data.conversationId}`);

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
                            this.error = 'Recipient not found';
                            break;
                        case 500:
                            console.error('Failed to get Coversation Internal Server Error:', error.response.data);
                            this.error = 'Failed to get Coversation Internal Server Error';
                            break;
                        default:
                            console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
                    }
                } else {
                        console.error('Error sending message:', error);
                        this.error = 'Network error';
                    }
            } finally {
                this.loading = false;
            }
        }
    }
}
</script>

<template>
    <div class="new-conversation p-4">
        <h2>Start New Conversation</h2>
        
        <div v-if="error" class="alert alert-danger">
            {{ error }}
        </div>

        <form @submit.prevent="startConversation">
            <div class="mb-3">
                <label class="form-label">Recipient Username:</label>
                <input 
                    type="text" 
                    class="form-control" 
                    v-model="recipientName" 
                    required
                    placeholder="Enter username"
                >
            </div>

            <div class="mb-3">
                <label class="form-label">Message:</label>
                <textarea 
                    class="form-control" 
                    v-model="message" 
                    required
                    rows="3"
                    placeholder="Type your message..."
                ></textarea>
            </div>

            <button 
                type="submit" 
                class="btn btn-primary" 
                :disabled="loading"
            >
                {{ loading ? 'Sending...' : 'Send Message' }}
            </button>
        </form>
    </div>
</template>

<style scoped>
.new-conversation {
    max-width: 600px;
    margin: 0 auto;
}
</style>