<template>
  <textarea
    ref="textarea"
    v-model="message"
    :rows="rows"
    class="resize-none"
    placeholder="Enterで改行してShift+Enterで送信します。"
    @keydown.delete="shrink"
    @keydown.enter.exact="grow"
    @keydown.enter.shift="send"
  ></textarea>
</template>

<script>
export default {
  props: {
    parent: {
      type: Object,
      default: null
    },
    rows: {
      type: Number,
      default: 1
    }
  },
  data() {
    return {
      message: ""
    };
  },
  methods: {
    grow() {
      //   ++this.rows;
      this.$emit("update:rows", this.rows + 1);
      const textarea = this.$refs.textarea;
      if (textarea.selectionStart !== textarea.selectionEnd) {
        const deletedText = this.message.substring(
          textarea.selectionStart,
          textarea.selectionEnd
        );
        // this.rows -= count(deletedText, "\n");
        this.$emit("update:rows", this.rows - count(deletedText, "\n"));
      }
    },
    send(e) {
      e.preventDefault();
      this.$emit("send", this.message, this.parent);
      this.$emit("update:rows", 1);
      this.message = "";
    },
    shrink() {
      const textarea = this.$refs.textarea;
      if (textarea.selectionStart === textarea.selectionEnd) {
        const char = this.message[textarea.selectionStart - 1];
        if (char === "\n") {
          //   --this.rows;
          this.$emit("update:rows", this.rows - 1);
        }
      } else {
        const deletedText = this.message.substring(
          textarea.selectionStart,
          textarea.selectionEnd
        );
        //this.rows -= count(deletedText, "\n")
        this.$emit("update:rows", this.rows - count(deletedText, "\n"));
      }
    }
  }
};
</script>
