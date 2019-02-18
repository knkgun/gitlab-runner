- [ ] Update the `.HelmChartMajor`, `.HelmChartMinor` and `.HelmChartPatch` to a specific release version
            git checkout -b update-runner-to-{{.Major}}-{{.Minor}}-0-rc1 && sed "s/^appVersion: .*/appVersion: {{.Major}}.{{.Minor}}.0-rc1/" Chart.yaml && git add Chart.yaml && git commit -m "Bump used Runner version to {{.Major}}.{{.Minor}}.0-rc1" -S && git push
        ## v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rc1 (TODAY_DATE_HERE)
    - [ ] add **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rc1** CHANGELOG entries and commit
        git add CHANGELOG.md && git commit -m "Update CHANGELOG for v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rc1" -S
    - [ ] bump version of the Helm Chart to `{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rc1`
        sed "s/^version: .*/version: {{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rc1/" Chart.yaml && git add Chart.yaml && git commit -m "Bump version to {{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rc1" -S && git push
    - [ ] tag and push **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rc1**:
        git tag -s v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rc1 -m "Version v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rc1" && git push origin v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rc1
    - [ ] create and push `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable` branch:
        git checkout -b {{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable && git push -u origin {{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable
    - [ ] checkout to `master`, bump version of the Helm Chart to `{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{inc .HelmChartPatch}}-beta` and push `master`:
        git checkout master; sed "s/^version: .*/version: {{.HelmChartMajor}}.{{.HelmChartMinor}}.{{inc .HelmChartPatch}}-beta/" Chart.yaml && git add Chart.yaml && git commit -m "Bump version to {{.HelmChartMajor}}.{{.HelmChartMinor}}.{{inc .HelmChartPatch}}-beta" -S && git push
    - [ ] check if Pipeline for `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable` is passing: [![pipeline status](https://gitlab.com/charts/gitlab-runner/badges/{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable/pipeline.svg)](https://gitlab.com/charts/gitlab-runner/commits/{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable)
        - [ ] add all required fixes to make `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable` Pipeline passing
    - [ ] `git checkout {{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable && git pull` in your local working copy!
            git checkout -b update-runner-to-{{.Major}}-{{.Minor}}-0-rcZ && sed "s/^appVersion: .*/appVersion: {{.Major}}.{{.Minor}}.0-rcZ/" Chart.yaml && git add Chart.yaml && git commit -m "Bump used Runner version to {{.Major}}.{{.Minor}}.0-rcZ" -S && git push
        - [ ] create Merge Request pointing `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable`: [link to MR here]
    - [ ] check if Pipeline for `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable` is passing: [![pipeline status](https://gitlab.com/charts/gitlab-runner/badges/{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable/pipeline.svg)](https://gitlab.com/charts/gitlab-runner/commits/{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable)
        - [ ] add all required fixes to make `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable` Pipeline passing
    - [ ] `git checkout {{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable && git pull` in your local working copy!
        ## v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ (TODAY_DATE_HERE)
    - [ ] add **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ** CHANGELOG entries and commit
        git add CHANGELOG.md && git commit -m "Update CHANGELOG for v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ" -S
    - [ ] bump version of the Helm Chart to `{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ`
        sed "s/^version: .*/version: {{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ/" Chart.yaml && git add Chart.yaml && git commit -m "Bump version to {{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ" -S && git push
    - [ ] tag and push **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ** and **{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable**:
        git tag -s v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ -m "Version v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ" && git push origin v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ {{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable
    - [ ] check if Pipeline for `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable` is passing: [![pipeline status](https://gitlab.com/charts/gitlab-runner/badges/{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable/pipeline.svg)](https://gitlab.com/charts/gitlab-runner/commits/{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable)
        - [ ] add all required fixes to make `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable` Pipeline passing
    - [ ] `git checkout {{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable && git pull` in your local working copy!
            git checkout -b update-runner-to-{{.Major}}-{{.Minor}}-0 && sed "s/^appVersion: .*/appVersion: {{.Major}}.{{.Minor}}.0/" Chart.yaml && git add Chart.yaml && git commit -m "Bump used Runner version to {{.Major}}.{{.Minor}}.0" -S && git push
        - [ ] create Merge Request pointing `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable`: [link to MR here]
    - [ ] check if Pipeline for `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable` is passing: [![pipeline status](https://gitlab.com/charts/gitlab-runner/badges/{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable/pipeline.svg)](https://gitlab.com/charts/gitlab-runner/commits/{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable)
        - [ ] add all required fixes to make `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable` Pipeline passing
    - [ ] `git checkout {{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable && git pull` in your local working copy!
        ## v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}} (TODAY_DATE_HERE)
    - [ ] add **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}** CHANGELOG entries and commit
        git add CHANGELOG.md && git commit -m "Update CHANGELOG for v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}" -S
    - [ ] bump version of the Helm Chart to `{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}`
        sed "s/^version: .*/version: {{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}/" Chart.yaml && git add Chart.yaml && git commit -m "Bump version to {{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}" -S && git push
    - [ ] tag and push **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}** and **{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable**:
        git tag -s v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}} -m "Version v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}" && git push origin v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}} {{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable
    - [ ] checkout to `master` and merge `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable` into `master` (only this one time, to update CHANGELOG.md and make the tag available for ./scripts/prepare-changelog-entries.rb in next stable release), push `master`:
        git checkout master; git merge --no-ff {{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable
    - [ ] update Runner's chart version [used by GitLab](https://gitlab.com/gitlab-org/gitlab-ce/blob/master/app/models/clusters/applications/runner.rb): [link to MR here]
    - [ ] check if Pipeline for `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable` is passing: [![pipeline status](https://gitlab.com/charts/gitlab-runner/badges/{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable/pipeline.svg)](https://gitlab.com/charts/gitlab-runner/commits/{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable)
        - [ ] add all required fixes to make `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable` Pipeline passing
    - [ ] `git checkout {{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable && git pull` in your local working copy!
            git checkout -b update-runner-to-{{.Major}}-{{.Minor}}-0-rcZ && sed "s/^appVersion: .*/appVersion: {{.Major}}.{{.Minor}}.0-rcZ/" Chart.yaml && git add Chart.yaml && git commit -m "Bump used Runner version to {{.Major}}.{{.Minor}}.0-rcZ" -S && git push
        - [ ] create Merge Request pointing `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable`: [link to MR here]
    - [ ] check if Pipeline for `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable` is passing: [![pipeline status](https://gitlab.com/charts/gitlab-runner/badges/{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable/pipeline.svg)](https://gitlab.com/charts/gitlab-runner/commits/{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable)
        - [ ] add all required fixes to make `{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable` Pipeline passing
    - [ ] `git checkout {{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable && git pull` in your local working copy!
        ## v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ (TODAY_DATE_HERE)
    - [ ] add **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ** CHANGELOG entries and commit
        git add CHANGELOG.md && git commit -m "Update CHANGELOG for v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ" -S
    - [ ] bump version of the Helm Chart to `{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ`
        sed "s/^version: .*/version: {{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ/" Chart.yaml && git add Chart.yaml && git commit -m "Bump version to {{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ" -S && git push
    - [ ] tag and push **v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ** and **{{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable**:
        git tag -s v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ -m "Version v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ" && git push origin v{{.HelmChartMajor}}.{{.HelmChartMinor}}.{{.HelmChartPatch}}-rcZ {{.HelmChartMajor}}-{{.HelmChartMinor}}-{{.HelmChartPatch}}-stable