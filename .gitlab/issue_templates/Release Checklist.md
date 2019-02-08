- [ ] Update the `.HelmChartMajor` and `.HelmChartMinor` to a specific release version
            git checkout -b update-runner-to-{{.Major}}-{{.Minor}}-0-rc1 && sed -i".bak" "s/^appVersion: .*/appVersion: {{.Major}}.{{.Minor}}.0-rc1/" Chart.yaml && rm Chart.yaml.bak && git add Chart.yaml && git commit -m "Bump used Runner version to {{.Major}}.{{.Minor}}.0-rc1" -S && git push -u origin update-runner-to-{{.Major}}-{{.Minor}}-0-rc1
        ## v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rc1 (TODAY_DATE_HERE)
    - [ ] add **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rc1** CHANGELOG entries and commit
        git add CHANGELOG.md && git commit -m "Update CHANGELOG for v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rc1" -S
    - [ ] bump version of the Helm Chart to `{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rc1`
        sed -i".bak" "s/^version: .*/version: {{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rc1/" Chart.yaml && rm Chart.yaml.bak && git add Chart.yaml && git commit -m "Bump version to {{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rc1" -S
    - [ ] tag and push **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rc1**:
        git tag -s v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rc1 -m "Version v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rc1" && git push origin v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rc1
    - [ ] create and push `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable` branch:
        git checkout -b {{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable && git push -u origin {{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable
    - [ ] checkout to `master`, bump version of the Helm Chart to `{{.HelmChartMajor}}.{{inc .HelmChartMinor}}.0-beta` and push `master`:
        git checkout master; sed -i".bak" "s/^version: .*/version: {{.HelmChartMajor}}.{{inc .HelmChartMinor}}.0-beta/" Chart.yaml && rm Chart.yaml.bak &&  git add Chart.yaml && git commit -m "Bump version to {{.HelmChartMajor}}.{{inc .HelmChartMinor}}.0-beta" -S && git push
    - [ ] check if Pipeline for `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable` is passing: [![pipeline status](https://gitlab.com/charts/gitlab-runner/badges/{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable/pipeline.svg)](https://gitlab.com/charts/gitlab-runner/commits/{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable)
        - [ ] add all required fixes to make `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable` Pipeline passing
    - [ ] `git checkout {{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable && git pull` in your local working copy!
            git checkout -b update-runner-to-{{.Major}}-{{.Minor}}-0-rcZ && sed -i".bak" "s/^appVersion: .*/appVersion: {{.Major}}.{{.Minor}}.0-rcZ/" Chart.yaml && rm Chart.yaml.bak && git add Chart.yaml && git commit -m "Bump used Runner version to {{.Major}}.{{.Minor}}.0-rcZ" -S && git push -u origin update-runner-to-{{.Major}}-{{.Minor}}-0-rcZ
        - [ ] create Merge Request pointing `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable`: [link to MR here]
    - [ ] check if Pipeline for `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable` is passing: [![pipeline status](https://gitlab.com/charts/gitlab-runner/badges/{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable/pipeline.svg)](https://gitlab.com/charts/gitlab-runner/commits/{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable)
        - [ ] add all required fixes to make `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable` Pipeline passing
    - [ ] `git checkout {{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable && git pull` in your local working copy!
        ## v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ (TODAY_DATE_HERE)
    - [ ] add **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ** CHANGELOG entries and commit
        git add CHANGELOG.md && git commit -m "Update CHANGELOG for v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ" -S
    - [ ] bump version of the Helm Chart to `{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ`
        sed -i".bak" "s/^version: .*/version: {{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ/" Chart.yaml && rm Chart.yaml.bak && git add Chart.yaml && git commit -m "Bump version to {{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ" -S
    - [ ] tag and push **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ** and **{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable**:
        git tag -s v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ -m "Version v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ" && git push origin v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ {{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable
    - [ ] check if Pipeline for `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable` is passing: [![pipeline status](https://gitlab.com/charts/gitlab-runner/badges/{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable/pipeline.svg)](https://gitlab.com/charts/gitlab-runner/commits/{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable)
        - [ ] add all required fixes to make `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable` Pipeline passing
    - [ ] `git checkout {{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable && git pull` in your local working copy!
            git checkout -b update-runner-to-{{.Major}}-{{.Minor}}-0 && sed -i".bak" "s/^appVersion: .*/appVersion: {{.Major}}.{{.Minor}}.0/" Chart.yaml && rm Chart.yaml.bak && git add Chart.yaml && git commit -m "Bump used Runner version to {{.Major}}.{{.Minor}}.0" -S && git push -u origin update-runner-to-{{.Major}}-{{.Minor}}-0
        - [ ] create Merge Request pointing `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable`: [link to MR here]
    - [ ] check if Pipeline for `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable` is passing: [![pipeline status](https://gitlab.com/charts/gitlab-runner/badges/{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable/pipeline.svg)](https://gitlab.com/charts/gitlab-runner/commits/{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable)
        - [ ] add all required fixes to make `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable` Pipeline passing
    - [ ] `git checkout {{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable && git pull` in your local working copy!
        ## v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0 (TODAY_DATE_HERE)
    - [ ] add **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0** CHANGELOG entries and commit
        git add CHANGELOG.md && git commit -m "Update CHANGELOG for v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0" -S
    - [ ] bump version of the Helm Chart to `{{.HelmChartMajor}}.{{.HelmChartMinor}}.0`
        sed -i".bak" "s/^version: .*/version: {{.HelmChartMajor}}.{{.HelmChartMinor}}.0/" Chart.yaml && rm Chart.yaml.bak && git add Chart.yaml && git commit -m "Bump version to {{.HelmChartMajor}}.{{.HelmChartMinor}}.0" -S
    - [ ] tag and push **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0** and **{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable**:
        git tag -s v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0 -m "Version v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0" && git push origin v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0 {{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable
    - [ ] checkout to `master` and merge `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable` into `master` (only this one time, to update CHANGELOG.md and make the tag available for `./scripts/prepare-changelog-entries.rb` in next stable release), push `master`:
        git checkout master; git merge --no-ff {{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable
- [ ] update Runner's chart version [used by GitLab](https://gitlab.com/gitlab-org/gitlab-ce/blob/master/app/models/clusters/applications/runner.rb): [link to MR here]

    - [ ] [create branch](https://gitlab.com/gitlab-org/gitlab-ce/branches/new?branch_name=update-gitlab-runner-helm-chart-to-{{.HelmChartMajor}}-{{.HelmChartMinor}}-0)
    - [ ] [create Merge Request](https://gitlab.com/gitlab-org/gitlab-ce/merge_requests/new?utf8=%E2%9C%93&merge_request[force_remove_source_branch]=1&merge_request[source_project_id]=13083&merge_request[target_project_id]=13083&merge_request[source_branch]=update-gitlab-runner-helm-chart-to-{{.HelmChartMajor}}-{{.HelmChartMinor}}-0&merge_request[target_branch]=master&merge_request[title]=Update+GitLab+Runner+Helm+Chart+to+{{.HelmChartMajor}}.{{.HelmChartMinor}}.0&merge_request[description]=%2Flabel%20~Verify%20~%22devops%3Averify%22%20~%22devops%3Aconfigure%22%20~%22dependency%20update%22%20%0A%2Fmilestone%20%25%22{{.Major}}.{{.Minor}}%22%0A%0A%23%23%20What%20does%20this%20MR%20do%3F%0A%0AUpdate%27s%20used%20GitLab%20Runner%20Helm%20Chart%20to%20version%20{{.HelmChartMajor}}.{{.HelmChartMinor}}.0%2C%20which%20uses%20GitLab%20Runner%20{{.Major}}.{{.Minor}}.0.%0A%0A%23%23%20What%20are%20the%20relevant%20issue%20numbers%3F%0A%0APart%20of%20gitlab-org%2Fgitlab-runner%23XXXX.%0A%0A%23%23%20Does%20this%20MR%20meet%20the%20acceptance%20criteria%3F%0A%0A-%20%5B%20%5D%20%5BChangelog%20entry%5D(https%3A%2F%2Fdocs.gitlab.com%2Fee%2Fdevelopment%2Fchangelog.html)%20added%2C%20if%20necessary%0A-%20%5B%20%5D%20%5BDocumentation%20created%2Fupdated%5D(https%3A%2F%2Fdocs.gitlab.com%2Fee%2Fdevelopment%2Fdocumentation%2Ffeature-change-workflow.html)%20via%20this%20MR%0A-%20%5B%20%5D%20Documentation%20reviewed%20by%20technical%20writer%20*or*%20follow-up%20review%20issue%20%5Bcreated%5D(https%3A%2F%2Fgitlab.com%2Fgitlab-org%2Fgitlab-ce%2Fissues%2Fnew%3Fissuable_template%3DDoc%2520Review)%0A-%20%5B%20%5D%20%5BTests%20added%20for%20this%20feature%2Fbug%5D(https%3A%2F%2Fdocs.gitlab.com%2Fee%2Fdevelopment%2Ftesting_guide%2Findex.html)%0A-%20%5B%20%5D%20Tested%20in%20%5Ball%20supported%20browsers%5D(https%3A%2F%2Fdocs.gitlab.com%2Fee%2Finstall%2Frequirements.html%23supported-web-browsers)%0A-%20%5B%20%5D%20Conforms%20to%20the%20%5Bcode%20review%20guidelines%5D(https%3A%2F%2Fdocs.gitlab.com%2Fee%2Fdevelopment%2Fcode_review.html)%0A-%20%5B%20%5D%20Conforms%20to%20the%20%5Bmerge%20request%20performance%20guidelines%5D(https%3A%2F%2Fdocs.gitlab.com%2Fee%2Fdevelopment%2Fmerge_request_performance_guidelines.html)%0A-%20%5B%20%5D%20Conforms%20to%20the%20%5Bstyle%20guides%5D(https%3A%2F%2Fgitlab.com%2Fgitlab-org%2Fgitlab-ee%2Fblob%2Fmaster%2FCONTRIBUTING.md%23style-guides)%0A-%20%5B%20%5D%20Conforms%20to%20the%20%5Bdatabase%20guides%5D(https%3A%2F%2Fdocs.gitlab.com%2Fee%2Fdevelopment%2FREADME.html%23database-guides)%0A-%20%5B%20%5D%20Link%20to%20e2e%20tests%20MR%20added%20if%20this%20MR%20has%20~%22Requires%20e2e%20tests%22%20label.%20See%20the%20%5BTest%20Planning%20Process%5D(https%3A%2F%2Fabout.gitlab.com%2Fhandbook%2Fengineering%2Fquality%2Ftest-engineering%2F).%0A-%20%5B%20%5D%20Security%20reports%20checked%2Fvalidated%20by%20reviewer)
    - [ ] Adjust and apply the patch to the GitLab CE sources

        <details>
        <summary>See the patch draft</summary>

        ```diff
        diff --git a/app/models/clusters/applications/runner.rb b/app/models/clusters/applications/runner.rb
        index 0c0247da1fb..f17da0bb7b1 100644
        --- a/app/models/clusters/applications/runner.rb
        +++ b/app/models/clusters/applications/runner.rb
        @@ -3,7 +3,7 @@
         module Clusters
           module Applications
             class Runner < ActiveRecord::Base
        -      VERSION = '{{.HelmChartMajor}}.{{dec .HelmChartMinor}}.0'.freeze
        +      VERSION = '{{.HelmChartMajor}}.{{.HelmChartMinor}}.0'.freeze

               self.table_name = 'clusters_applications_runners'

        diff --git a/changelogs/unreleased/update-gitlab-runner-helm-chart-to-{{.HelmChartMajor}}-{{.HelmChartMinor}}-0.yml b/changelogs/unreleased/update-gitlab-runner-helm-chart-to-{{.HelmChartMajor}}-{{.HelmChartMinor}}-0.yml
        new file mode 100644
        index 00000000000..7d92929221f
        --- /dev/null
        +++ b/changelogs/unreleased/update-gitlab-runner-helm-chart-to-{{.HelmChartMajor}}-{{.HelmChartMinor}}-0.yml
        @@ -0,0 +1,5 @@
        +---
        +title: Update GitLab Runner Helm Chart to {{.HelmChartMajor}}.{{.HelmChartMinor}}.0
        +merge_request: XXX
        +author:
        +type: other
        diff --git a/spec/models/clusters/applications/runner_spec.rb b/spec/models/clusters/applications/runner_spec.rb
        index 3d0735c6d0b..8ad41e997c2 100644
        --- a/spec/models/clusters/applications/runner_spec.rb
        +++ b/spec/models/clusters/applications/runner_spec.rb
        @@ -46,7 +46,7 @@ describe Clusters::Applications::Runner do
             it 'should be initialized with 4 arguments' do
               expect(subject.name).to eq('runner')
               expect(subject.chart).to eq('runner/gitlab-runner')
        -      expect(subject.version).to eq('{{.HelmChartMajor}}.{{dec .HelmChartMinor}}.0')
        +      expect(subject.version).to eq('{{.HelmChartMajor}}.{{.HelmChartMinor}}.0')
               expect(subject).to be_rbac
               expect(subject.repository).to eq('https://charts.gitlab.io')
               expect(subject.files).to eq(gitlab_runner.files)
        @@ -64,7 +64,7 @@ describe Clusters::Applications::Runner do
               let(:gitlab_runner) { create(:clusters_applications_runner, :errored, runner: ci_runner, version: '0.1.13') }

               it 'should be initialized with the locked version' do
        -        expect(subject.version).to eq('{{.HelmChartMajor}}.{{dec .HelmChartMinor}}.0')
        +        expect(subject.version).to eq('{{.HelmChartMajor}}.{{.HelmChartMinor}}.0')
               end
             end
           end
        ```

        Adjust the patch (set proper current version; set the Merge Request ID in the CHANGELOG entry file) and save
        the patch at `/tmp/patch.gitlab-ce`.

        Go to your local GitLab CE working directory:

        ```bash
        git pull && \
            git checkout update-gitlab-runner-helm-chart-to-{{.HelmChartMajor}}-{{.HelmChartMinor}}-0 && \
            git apply /tmp/patch.gitlab-ce && \
            git add . && \
            git commit -m "Update GitLab Runner Helm Chart to {{.HelmChartMajor}}.{{.HelmChartMinor}}.0" && \
            git push
        ```

        </details>

- [ ] update Runner's chart version [used by GitLab chart](https://gitlab.com/charts/gitlab/blob/master/requirements.yaml#L16): [link to MR here]

    - [ ] [create branch](https://gitlab.com/charts/gitlab/branches/new?branch_name=update-gitlab-runner-helm-chart-to-{{.HelmChartMajor}}-{{.HelmChartMinor}}-0)
    - [ ] [create Merge Request](https://gitlab.com/charts/gitlab/merge_requests/new?utf8=%E2%9C%93&merge_request[force_remove_source_branch]=1&merge_request[source_project_id]=3828396&merge_request[target_project_id]=3828396&merge_request[source_branch]=update-gitlab-runner-helm-chart-to-{{.HelmChartMajor}}-{{.HelmChartMinor}}-0&merge_request[target_branch]=master&merge_request[title]=Update%20GitLab%20Runner%20Helm%20Chart%20to%20{{.HelmChartMajor}}.{{.HelmChartMinor}}.0%2F{{.Major}}.{{.Minor}}.0&merge_request[description]=%2Fmilestone%20%25%22{{.Major}}.{{.Minor}}%22%0A%0AUpdate%20the%20gitlab%2Fgitlab-runner%20requirement%20for%20v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0%20containing%20GitLab%20Runner%20{{.Major}}.{{.Minor}}.0)
    - [ ] Adjust and apply the patch to the GitLab Helm Chart sources

        <details>
        <summary>See the patch draft</summary>

        ```diff
        diff --git a/changelogs/unreleased/update-gitlab-runner-to-{{.HelmChartMajor}}-{{.HelmChartMinor}}-0.yml b/changelogs/unreleased/update-gitlab-runner-to-{{.HelmChartMajor}}-{{.HelmChartMinor}}-0.yml
        new file mode 100644
        index 00000000..5fd73668
        --- /dev/null
        +++ b/changelogs/unreleased/update-gitlab-runner-to-{{.HelmChartMajor}}-{{.HelmChartMinor}}-0.yml
        @@ -0,0 +1,5 @@
        +---
        +title: Update gitlab-runner to {{.HelmChartMajor}}.{{.HelmChartMinor}}.0/{{.Major}}.{{.Minor}}.0
        +merge_request: XXX
        +author:
        +type: other
        diff --git a/requirements.yaml b/requirements.yaml
        index 82a0192a..c6293a8d 100644
        --- a/requirements.yaml
        +++ b/requirements.yaml
        @@ -13,6 +13,6 @@ dependencies:
           repository: https://kubernetes-charts.storage.googleapis.com/
           condition: postgresql.install
         - name: gitlab-runner
        -  version: {{.HelmChartMajor}}.{{dec .HelmChartMinor}}.0
        +  version: {{.HelmChartMajor}}.{{.HelmChartMinor}}.0
           repository: https://charts.gitlab.io/
           condition: gitlab-runner.install
        ```

        Adjust the patch (set proper current version; set the Merge Request ID in the CHANGELOG entry file) and save
        the patch at `/tmp/patch.gitlab-helm-chart`.

        Go to your local GitLab Helm Chart working directory:

        ```bash
        git pull && \
            git checkout update-gitlab-runner-helm-chart-to-{{.HelmChartMajor}}-{{.HelmChartMinor}}-0 && \
            git apply /tmp/patch.gitlab-helm-chart && \
            git add . && \
            git commit -m "Update GitLab Runner Helm Chart to {{.HelmChartMajor}}.{{.HelmChartMinor}}.0" && \
            git push
        ```

        </details>
    - [ ] check if Pipeline for `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable` is passing: [![pipeline status](https://gitlab.com/charts/gitlab-runner/badges/{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable/pipeline.svg)](https://gitlab.com/charts/gitlab-runner/commits/{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable)
        - [ ] add all required fixes to make `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable` Pipeline passing
    - [ ] `git checkout {{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable && git pull` in your local working copy!
            git checkout -b update-runner-to-{{.Major}}-{{.Minor}}-0-rcZ && sed -i".bak" "s/^appVersion: .*/appVersion: {{.Major}}.{{.Minor}}.0-rcZ/" Chart.yaml && rm Chart.yaml.bak && git add Chart.yaml && git commit -m "Bump used Runner version to {{.Major}}.{{.Minor}}.0-rcZ" -S && git push -u origin update-runner-to-{{.Major}}-{{.Minor}}-0-rcZ
        - [ ] create Merge Request pointing `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable`: [link to MR here]
    - [ ] check if Pipeline for `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable` is passing: [![pipeline status](https://gitlab.com/charts/gitlab-runner/badges/{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable/pipeline.svg)](https://gitlab.com/charts/gitlab-runner/commits/{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable)
        - [ ] add all required fixes to make `{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable` Pipeline passing
    - [ ] `git checkout {{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable && git pull` in your local working copy!
        ## v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ (TODAY_DATE_HERE)
    - [ ] add **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ** CHANGELOG entries and commit
        git add CHANGELOG.md && git commit -m "Update CHANGELOG for v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ" -S
    - [ ] bump version of the Helm Chart to `{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ`
        sed -i".bak" "s/^version: .*/version: {{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ/" Chart.yaml && rm Chart.yaml.bak && git add Chart.yaml && git commit -m "Bump version to {{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ" -S
    - [ ] tag and push **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ** and **{{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable**:
        git tag -s v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ -m "Version v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ" && git push origin v{{.HelmChartMajor}}.{{.HelmChartMinor}}.0-rcZ {{.HelmChartMajor}}-{{.HelmChartMinor}}-0-stable