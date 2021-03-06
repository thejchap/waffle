(($) => {
  /**
   * Main chat client. Maintains state and also handles UI
   */
  class Client {
    constructor() {
      this.ENDPOINTS = {
        SSE: '/sse',
        MESSAGES: '/api/messages'
      };

      this.state = {
        messageIds: new Set(),
        messages: []
      };
    }

    /**
     * Initializes the client, assigns an ID, sets up SSE listener and event
     * bindings, and loads the initial list of messages
     */
    async run() {
      this.id = await this._fetchClientId();

      this._setupEventSource();
      this._loadMessages();
      this._setupActions();

      return this.id;
    }

    /**
     * Setup SSE event source. At this point only handle 1 type of message
     */
    _setupEventSource() {
      const { ENDPOINTS: { SSE } } = this;
      const eventSource = new EventSource(SSE);

      eventSource.onmessage = ({ data }) => {
        const msg = JSON.parse(data);

        if (msg.keepalive === true) {
          return;
        }

        if (msg.sender === this.id) {
          return;
        }

        this._pushMessage(JSON.parse(data));
        this._render();
      };

      this.eventSource = eventSource;
    }

    /**
     * UI Event handlers
     */
    _setupActions() {
      $('form').submit((e) => {
        e.preventDefault();
        const content = $('#content').val();
        $('#content').val('');
        this._createMessage(content);
      });

      $('#send').removeClass('disabled');
    }

    /**
     * Push a newly created message into the local store. If this succeeds, post
     * to the server
     */
    async _createMessage(content) {
      const timestamp = new Date().getTime();
      const id = await this._genId();
      const { id: sender, ENDPOINTS: { MESSAGES: url } } = this;
      const message = { content, id, timestamp, sender };

      if (!this._pushMessage(message)) {
        return false;
      }

      this._render();

      return $.ajax({
        type: 'POST',
        url,
        dataType: 'json',
        contentType: 'application/json',
        data: JSON.stringify(message)
      });
    }

    /**
     * Essentially a session ID. If it exists in localStorage, grab from there.
     * otherwise generate a new one an save in localStorage
     */
    async _fetchClientId() {
      const key = 'waffle:clientId';
      const existing = localStorage.getItem(key);

      if (existing) {
        return existing;
      }

      const id = await this._genId();

      localStorage.setItem(key, id);

      return id;
    }

    /**
     * Generates a unique ID. Combines the current time with a random string and
     * hashes using SHA-256. Return the first 9 characters of the hex string
     * representation of the hash.
     * Possible IDs: 16^9 ~= 68 billion
     */
    async _genId() {
      const encoder = new TextEncoder();
      const now = new Date().getTime();
      const rand = Math.random().toString(36).substr(2, 9);
      const message = `${now}:${rand}`;
      const data = encoder.encode(message);
      const buffer = await crypto.subtle.digest('SHA-256', data);
      const arr = new Uint8Array(buffer);

      return [...arr].map((val) => {
        return val.toString(16).padStart(2, '0');
      }).join('').substr(0, 9);
    }

    /**
     * Push a message into the local store unless it exists already
     */
    _pushMessage(message) {
      const { state: { messageIds, messages } } = this;

      if (messageIds.has(message.id)) {
        return false;
      }

      messageIds.add(message.id);
      messages.push(message);

      return true;
    }

    /**
     * Fetch the initial list of messages from the server
     */
    _loadMessages() {
      const { ENDPOINTS: { MESSAGES: url } } = this;

      $.ajax({ type: 'GET', url, dataType: 'json', }).then((messages) => {
        messages.forEach(this._pushMessage.bind(this));
        this._render();
      });
    }

    /**
     * Render the message list UI
     */
    _render() {
      const { state: { messages } } = this;
      const $list = $('#messages');

      $list.empty();

      messages.sort((a, b) => a.timestamp - b.timestamp).forEach((msg, i) => {
        const { content, sender } = msg;
        const $msg = $('<li />');
        const $content = $('<div />');
        const $client = $('<small />');

        $client
          .text(sender)
          .addClass('d-block px-2 mb-1 mt-2 client');

        $content
          .text(content)
          .addClass('d-inline-block px-3 py-2 rounded-pill');

        $msg
          .addClass('list-group-item py-1')
          .append($client)
          .append($content);

        // Visually differentiate between my messages and other people's messages
        if (msg.sender === this.id) {
          $msg.addClass('text-right');
          $content.addClass('bg-primary text-white');
        } else {
          $content.addClass('bg-light');
        }

        // Group together messages from the same sender
        if (msg.sender === this.id || (i > 0 && msg.sender === messages[i - 1].sender)) {
          $client.addClass('d-none').removeClass('d-block');
        }

        $list.append($msg);
      });

      $('.messages-wrapper').scrollTop($list.prop('scrollHeight'));
    }
  }

  // Bootstrap the app
  $(() => new Client().run());
})(window.jQuery);
