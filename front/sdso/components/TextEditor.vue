<template>
  <div ref="editor" id="editor" class="w-100 h-100 mt-5"></div>
</template>

<script>
let editor = null;
export default {
  props: {
    fontSize: {
      type: Number,
      default: 20
    },
    mode: {
      type: String,
      required: true
    },
    readOnly: {
      type: Boolean,
      default: true
    },
    theme: {
      type: String,
      default: "ace/theme/xcode"
    },
    value: {
      type: String,
      default: ""
    }
  },
  mounted() {
    this.setupAce();
  },
  watch: {
    mode: {
      immediate: true,
      handler(newMode) {
        if (!editor) {
          return;
        }
        editor.getSession().setMode(newMode);
      }
    },
    theme: {
      immediate: true,
      handler(newTheme) {
        if (!editor) {
          return;
        }
        editor.setTheme(newTheme);
      }
    },
    value: {
      immediate: true,
      handler(newValue) {
        if (!editor) {
          return;
        }
        editor.setValue(newValue);
      }
    }
  },
  methods: {
    setupAce() {
      editor = ace.edit("editor");
      editor.setTheme(this.theme);
      editor.setFontSize(this.fontSize);
      editor.setReadOnly(this.readOnly);
      editor.getSession().setMode(this.mode);
      editor.setValue(this.value);
      editor.getSession().on("change", this.updateEditorHeight);
    },
    updateEditorHeight() {
      if (!this.$refs.editor) {
        return;
      }
      const newHeight =
        editor.getSession().getScreenLength() * editor.renderer.lineHeight +
        editor.renderer.scrollBar.getWidth();
      this.$refs.editor.style.height = newHeight.toString() + "px";
      editor.resize();
    }
  }
};
</script>


<style>
.ace_scrollbar-v {
  display: none;
}
</style>