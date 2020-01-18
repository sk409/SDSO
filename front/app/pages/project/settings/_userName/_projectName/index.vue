<template>
  <div class="container">
    <div class="row">
      <div class="col-12 border-bottom h3">ブランチ</div>
    </div>
    <div class="row mt-3">
      <div class="col-12">
        <el-card>
          <div slot="header">保護</div>
          <div>
            <div
              v-for="branchProtectionRule in branchProtectionRules"
              :key="branchProtectionRule.ID"
            >
              {{ branchProtectionRule.BranchName }}
            </div>
            <div class="text-center mt-3">
              <el-button
                type="primary"
                @click="showAddingBranchProtectionRuleDialog"
                >追加</el-button
              >
            </div>
          </div>
        </el-card>
      </div>
    </div>
    <el-dialog :visible.sync="addingBranchProtectionRuleDialog.visible">
      <el-form
        ref="branchProtectionRuleForm"
        :model="branchProtectionRuleForm"
        :rules="branchProtectionRuleFormRules"
      >
        <el-form-item label="ブランチ名" prop="branchName">
          <el-input v-model="branchProtectionRuleForm.branchName"></el-input>
        </el-form-item>
        <div class="text-center">
          <el-button type="primary" @click="addBranchProtectionRule"
            >追加</el-button
          >
        </div>
      </el-form>
    </el-dialog>
  </div>
</template>

<script>
let user = null;
export default {
  layout: "Project",
  data() {
    const that = this;
    const duplicateNameValidator = function(rule, value, callback) {
      for (const branchProtectionRule of that.branchProtectionRules) {
        if (value === branchProtectionRule.BranchName) {
          callback(new Error("同じ名前のブランチが既に存在しています"));
          return;
        }
      }
      callback();
    };
    return {
      addingBranchProtectionRuleDialog: {
        visible: false
      },
      branchProtectionRuleForm: {
        branchName: ""
      },
      branchProtectionRuleFormRules: {
        branchName: [
          {
            required: true,
            message: "ブランチ名を入力してください",
            trigger: "blur"
          },
          { validator: duplicateNameValidator, trigger: "blur" }
        ]
      },
      project: null
    };
  },
  computed: {
    branchProtectionRules() {
      return this.project ? this.project.branchProtectionRules : [];
    },
    pathParamUserName() {
      return this.$route.params.userName;
    },
    pathParamProjectName() {
      return this.$route.params.projectName
        ? this.$route.params.projectName
        : this.$route.params.pathMatch;
    }
  },
  created() {
    this.$ajax.get(this.$urls.user, {}, { withCredentials: true }, response => {
      if (response.status !== 200) {
        return;
      }
      user = response.data;
      const data = {
        name: this.pathParamProjectName,
        user_id: user.ID
      };
      this.$ajax.get(this.$urls.projects, data, {}, response => {
        if (response.status !== 200) {
          return;
        }
        if (response.data.length !== 1) {
          return;
        }
        this.project = response.data[0];
        const data = {
          project_id: this.project.ID
        };
        this.$ajax.get(this.$urls.branchProtectionRules, data, {}, response => {
          this.$set(this.project, "branchProtectionRules", response.data);
        });
      });
    });
  },
  methods: {
    addBranchProtectionRule() {
      this.$refs.branchProtectionRuleForm.validate(valid => {
        if (!valid) {
          return;
        }
        const data = {
          BranchName: this.branchProtectionRuleForm.branchName,
          ProjectID: this.project.ID
        };
        this.$ajax.post(
          this.$urls.branchProtectionRules,
          data,
          {},
          response => {
            if (response.status !== 200) {
              return;
            }
            this.project.branchProtectionRules.push(response.data);
            this.addingBranchProtectionRuleDialog.visible = false;
          }
        );
      });
    },
    showAddingBranchProtectionRuleDialog() {
      this.addingBranchProtectionRuleDialog.visible = true;
    }
  }
};
</script>
