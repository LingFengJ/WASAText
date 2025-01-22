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
                console.error('Error:', error);
                this.error = 'Network error';
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