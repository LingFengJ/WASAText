<script>
import ConversationList from '@/components/ConversationList.vue'

export default {
  name: 'ConversationsView',
  components: {
    ConversationList
  },
  data() {
    return {
      searchQuery: '',
      searchResults: [],
      loading: false,
      showResults: false,
      error: null,
    }
  },
  created() {
    const username = this.$route.query.username;
    if (username) {
        this.recipientName = username;
    }
    // Add click outside listener when component is created
    document.addEventListener('click', this.handleClickOutside);
  },
  beforeDestroy() {
      // Clean up the listener when component is destroyed
      document.removeEventListener('click', this.handleClickOutside);
  },
  methods: {
    async performSearch() {
      this.loading = true;
      try {
        console.log('Performing search for:', this.searchQuery);
        const response = await this.$axios.get(
          `/users/search?query=${encodeURIComponent(this.searchQuery)}`,
          {
            headers: {
              'Authorization': `Bearer ${sessionStorage.getItem('authToken')}`
            }
          }
        );

        console.log('Search response status:', response.status);
        if (response.status === 200) {
          console.log('Search results:', response.data);
          this.searchResults = response.data;
          this.showResults = true;
        } else {
          console.error('Search failed:', response.status, response.data);
          this.searchResults = [];
          this.showResults = false;
        }
      } catch (error) {
        console.error('Search error:', error);
        this.searchResults = [];
        this.showResults = false;
      } finally {
        this.loading = false;
      }
    },
    async startChat(user) {
      if (user.username === sessionStorage.getItem('username')) {
        this.error = "You cannot start a chat with yourself";
        setTimeout(() => this.error = null, 3000); // Clear error after 3 seconds
        this.searchQuery = '';
        this.showResults = false;
        return;
      }
      try {
          // First get all existing conversations
          const convResponse = await this.$axios.get('/conversations', {
              headers: {
                  'Authorization': `Bearer ${sessionStorage.getItem('authToken')}`
              }
          });

          // Look for existing individual chat with this user
          const existingChat = convResponse.data.find(conv => 
              conv.type === 'individual' && conv.name === user.username
          );

          if (existingChat) {
              // If chat exists, just navigate to it
              this.$router.push(`/conversations/${existingChat.id}`);
          } else {
              // If no chat exists, redirect to NewConversation with pre-filled username
              this.$router.push({
                  path: '/new-conversation',
                  query: { username: user.username }
              });
          }

          // Clear search
          this.searchQuery = '';
          this.showResults = false;

      } catch (error) {
          console.error('Error checking conversations:', error);
      }
  },

    handleClickOutside(event) {
        if (!event.target.closest('.search-container')) {
            this.showResults = false;
        }
    },
  }
}
</script>

<template>
  <div class="conversations-view">
    
    <div class="d-flex justify-content-between align-items-center mb-4">
      <!-- Add error alert -->
        <div v-if="error" class="alert alert-danger position-absolute top-30 start-50 translate-middle-x mt-3">
          {{ error }}
        </div>
      <h1>Conversations</h1>
      <div class="search-container">
        <div class="input-group">
          <input
            type="text"
            v-model="searchQuery"
            class="form-control"
            placeholder="Search users..."
            @input="performSearch"
            @focus="() => { this.searchQuery = ''; performSearch(); }"
          />
          <span class="input-group-text">
            <i class="bi bi-search"></i>
          </span>
        </div>
        
        <!-- Search Results Dropdown -->
        <div v-if="showResults" class="search-results">
          <div v-if="loading" class="p-2 text-center">
            <div class="spinner-border spinner-border-sm" role="status">
              <span class="visually-hidden">Loading...</span>
            </div>
          </div>
          <div v-else-if="searchResults && searchResults.length > 0" class="list-group">
            <button
              v-for="user in searchResults"
              :key="user.id"
              class="list-group-item list-group-item-action"
              @click="startChat(user)"
            >
              <i class="bi bi-person-circle me-2"></i>
              {{ user.username }}
            </button>
          </div>
          <div v-else class="p-2 text-center text-muted">
            No users found
          </div>
        </div>
      </div>
    </div>

    <ConversationList />
  </div>
</template>

<style scoped>
.conversations-view {
  padding: 20px;
}

.search-container {
  position: relative;
  width: 300px;
}

.search-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  z-index: 1000;
  max-height: 300px;
  overflow-y: auto;
}

.list-group-item {
  cursor: pointer;
  border: none;
  border-bottom: 1px solid #eee;
}

.list-group-item:hover {
  background-color: #f8f9fa;
}

.bi-search {
  color: #6c757d;
}
</style>