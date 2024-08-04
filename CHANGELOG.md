# Changelog

---
## [0.0.2-match-service](https://github.com/hedon954/go-matcher/compare/v0.0.1-init..v0.0.2-match-service) - 2024-08-04

finish match service layer

### ⚙️ Miscellaneous Chores

- **(ci)** do not fail Github Action ci when test coverage is below 50. - ([03b1be0](https://github.com/hedon954/go-matcher/commit/03b1be0c71a4b1f09ec5d6aa0dbca018c6da8646)) - wangjiahan
- **(ci)** change git cliff config - ([33dd949](https://github.com/hedon954/go-matcher/commit/33dd949c42bb7d25a3dff8b2e1766c506a49ab92)) - hedon954
- **(ci)** change git cliff config - ([834db10](https://github.com/hedon954/go-matcher/commit/834db10d2c1e7f06f27721c77412f6d853756bd3)) - hedon954
- **(ci)** fix github action config - ([87b12b2](https://github.com/hedon954/go-matcher/commit/87b12b2f323bc6b06cf53d3c4d086b53602c4045)) - hedon954
- **(lint)** replace `gomnd` with `mnd` by golangci-lint prompt. - ([2bd894f](https://github.com/hedon954/go-matcher/commit/2bd894f045e873d2af58ed45a3a3566affd3a5e3)) - wangjiahan
- add comments to matcher interface - ([0189dce](https://github.com/hedon954/go-matcher/commit/0189dce945fc21b895188ee266004bfb63542ea4)) - wangjiahan
- optimize matcher implement code order - ([207aba7](https://github.com/hedon954/go-matcher/commit/207aba7869a86261cb17aa6e7330aab47dbe2ee5)) - wangjiahan
- add comment to group interface - ([7e3107e](https://github.com/hedon954/go-matcher/commit/7e3107e6991df7c3fd584e048453bb2e91670d90)) - wangjiahan
- delete `ReadyToMatch` interface - ([a9193c9](https://github.com/hedon954/go-matcher/commit/a9193c9e37d7fe9be499141390b9e42379654229)) - wangjiahan

### ⛰️ Features

- **(ci)** replace golangci-lint with Github Action and add more hooks in pre commit. - ([20de6a2](https://github.com/hedon954/go-matcher/commit/20de6a2f90b21168fa72ce92b7410d595a36574b)) - wangjiahan
- **(ci)** run test with coverage and send coverage result to `codecov`. - ([0ed649c](https://github.com/hedon954/go-matcher/commit/0ed649c1cf9543af688c0298959a7ac4f607262e)) - wangjiahan
- **(ci)** add race detector to pre commit hook and Github Action workflow. - ([7ad027c](https://github.com/hedon954/go-matcher/commit/7ad027c995fc3de566b057832257081e6dab45f2)) - wangjiahan
- **(lint)** update golangci-lint config. - ([79ab912](https://github.com/hedon954/go-matcher/commit/79ab912202a5e268183c054cb602ea42ec16d8df)) - wangjiahan
- **(service)** finish `SetNearbyJoinGroup` and `SetRecentJoinGroup` apis - ([c79a827](https://github.com/hedon954/go-matcher/commit/c79a82708ffc453a8eccadb1d918346cc920e283)) - wangjiahan
- **(service)** finish `invite` api - ([095a507](https://github.com/hedon954/go-matcher/commit/095a5079453aded276725dddbe53e7987146b41d)) - wangjiahan
- **(service)** add `matchChannel` and `closeChannel` to service and refactor the unit tests - ([558b70a](https://github.com/hedon954/go-matcher/commit/558b70a09d3d75c6491afdc0b3922bd97c6f00dd)) - hedon954
- add example for testing pre-commit and Github Action. - ([8aa0016](https://github.com/hedon954/go-matcher/commit/8aa001696e371ba8e795077185f30d66c84926f9)) - hedon954
- build go matcher skeleton. - ([628385b](https://github.com/hedon954/go-matcher/commit/628385bfd2c9dd39813cc130215facea4a88995e)) - hedon954
- add `playerManager` and move basic interfaces to `common` package. - ([b8e137c](https://github.com/hedon954/go-matcher/commit/b8e137cfbe242fa9a08f56bece586abd0eb493cd)) - wangjiahan
- finish CreateGroup, InviteFriend, HandleInvite, Kick and ExitGroup implementaion. - ([893f195](https://github.com/hedon954/go-matcher/commit/893f1955a64017f04d13b030eec6a5bdcca54b00)) - wangjiahan
- support part of match api - ([bf148e3](https://github.com/hedon954/go-matcher/commit/bf148e33f7ff477da2779fc5f4574377f0e8fb0a)) - wangjiahan
- finish `CreateGroup`, `DissolveGroup`, `EneterGroup` and `ExitGroup` matcher apis - ([c3cd952](https://github.com/hedon954/go-matcher/commit/c3cd9520f41457ceb51ae8f48f2e131c76c5dcf8)) - wangjiahan
- finish `kick_player` service - ([379476b](https://github.com/hedon954/go-matcher/commit/379476b2659067c6834820ac1afd92c6034b04e6)) - wangjiahan
- finish `RefuseInvite` api - ([a1806ff](https://github.com/hedon954/go-matcher/commit/a1806ff9753ea170dca7d41ccf7767ef7989c1e2)) - hedon954
- finish `AcceptInvite` api - ([1f5ba09](https://github.com/hedon954/go-matcher/commit/1f5ba09152675082bbac2c675ccb1f95bf322bc7)) - hedon954
- add and finish `SetVoiceState` api - ([60748c5](https://github.com/hedon954/go-matcher/commit/60748c56e89d5eac9d2584c9eed89563ffdd0c29)) - hedon954
- finish the basic logic of `StartMatch` and `CancelMatch` - ([9d7b154](https://github.com/hedon954/go-matcher/commit/9d7b1548a3ac13f7800b40b9b89f4c6a8866e2e9)) - hedon954
- send group to match queue when start match - ([e7e0570](https://github.com/hedon954/go-matcher/commit/e7e05705a7c1e1e6a226a66c0eaa44aac48a51e7)) - hedon954

### 🐛 Bug Fixes

- **(ci)** add `fi` for check test coverage if statement check. - ([6c42db1](https://github.com/hedon954/go-matcher/commit/6c42db1ae59b4224af0ebaef2eaf0392f8a96446)) - wangjiahan
- **(ci)** just lint changed files - ([2e9863b](https://github.com/hedon954/go-matcher/commit/2e9863bac062381a6f661a980daa3b1bce775d45)) - wangjiahan
- **(ci)** fix git cliff action version to v3 - ([47444b0](https://github.com/hedon954/go-matcher/commit/47444b0b48856f0d6f55dcb891340d67983ee81c)) - hedon954
- **(ci)** cat NEW_CHANGELOG - ([60ce60f](https://github.com/hedon954/go-matcher/commit/60ce60f52e1fae0c8c956670f48f99eb5bee4a7b)) - hedon954
- **(ci)** fix git cliff args - ([3b7ee92](https://github.com/hedon954/go-matcher/commit/3b7ee923fe91630ce121ee5c30c59ee406a3b48e)) - hedon954
- **(lint)** optimize code with golangci-lint recommends. - ([17bcd2c](https://github.com/hedon954/go-matcher/commit/17bcd2cc7cb1825e2cd1b74af117c862588e6387)) - wangjiahan
- **(lint)** optimize code according to golangci-lint prompts - ([2f14f51](https://github.com/hedon954/go-matcher/commit/2f14f517ab6a331486cb4663d1f0436f25ff39a1)) - wangjiahan
- **(lint)** optimize code according to golangci-lint prompts - ([df055dd](https://github.com/hedon954/go-matcher/commit/df055ddcc2949b4048579d4d3ea983bd98afb15b)) - wangjiahan
- **(service)** push group state when cancel match - ([f1226cb](https://github.com/hedon954/go-matcher/commit/f1226cbae69796de75021934c8342d1839e561d6)) - hedon954

### 📚 Documentation

- **(changelog)** fix `Github` to `GitHub`. - ([08d280e](https://github.com/hedon954/go-matcher/commit/08d280e9a0005a9ec1aeb2fd9686d40246326095)) - hedon954
- **(changelog)** reset changelog - ([2ab674a](https://github.com/hedon954/go-matcher/commit/2ab674a6c53d72be318cd22bbc82aa1a7386e2fd)) - hedon954
- **(changelog)** reset changelog - ([d5339d7](https://github.com/hedon954/go-matcher/commit/d5339d748dc11f5f589868d952c367c72f1cef94)) - hedon954
- **(changelog)** reset changelog - ([f9ff20c](https://github.com/hedon954/go-matcher/commit/f9ff20c5398bdb6374cbd76c867a871b8c0eada1)) - hedon954

### 🚜 Refactor

- **(service)** change `HandoverCapation` to `ChangeRole` to support more cases - ([5b0ce2e](https://github.com/hedon954/go-matcher/commit/5b0ce2e15ecdb4e6ddd567b0c4258f4e91880310)) - wangjiahan
- delete useless files - ([d1a385e](https://github.com/hedon954/go-matcher/commit/d1a385ea5b98350b33b08bc2828baed8dfb9eb64)) - wangjiahan
- rename package `manager` to `repository` - ([96f80c9](https://github.com/hedon954/go-matcher/commit/96f80c9cc7a8e36041c4b91e7e15dc5dcd70d0d0)) - wangjiahan
- rename package `matcher` to `service` - ([35c80e6](https://github.com/hedon954/go-matcher/commit/35c80e6a805ecd93dd5648e7c98afc6e472e50cb)) - wangjiahan
- change `InvitationSrcType` to `EnterGroupSourceType` for better understand - ([f4ce789](https://github.com/hedon954/go-matcher/commit/f4ce789d738a74e734268216bceef70f13788be4)) - wangjiahan
- redesign the `CanPlayTogether` interface and fix the unit tests of the `EnterGroup` api - ([dbff935](https://github.com/hedon954/go-matcher/commit/dbff9353cc90373b1158626e62ffa9f75aa57782)) - hedon954
- rename `closeChannel` to `stopChannel` - ([1c960e8](https://github.com/hedon954/go-matcher/commit/1c960e82f04284338cebac16a0181a8982f51f18)) - hedon954
- remove handle system close logic from service layer - ([786c104](https://github.com/hedon954/go-matcher/commit/786c1041710df357f5495229049b28ea44b0104c)) - hedon954

### 🧪 Tests

- **(matcher)** add more unit tests for macther default implement - ([5904d19](https://github.com/hedon954/go-matcher/commit/5904d196fdecc16b3390c3bc3d5abc5d34d06fb8)) - wangjiahan
- **(matcher)** check user state when create group - ([47913e0](https://github.com/hedon954/go-matcher/commit/47913e03426b1122552941897a7fcd35ec0c2b1a)) - wangjiahan
- **(matcher)** add `HandoverCaptain` unit tests - ([d69906e](https://github.com/hedon954/go-matcher/commit/d69906eded8fcb6a9700963ade96938cd65e815a)) - wangjiahan
- add more unit test cases to improve test coverage. - ([780b66b](https://github.com/hedon954/go-matcher/commit/780b66b21be9cc7c4bb40c018129fd3c090aecc0)) - wangjiahan
- add more unit test cases to improve coverage. - ([79bf6a9](https://github.com/hedon954/go-matcher/commit/79bf6a9e3f5f8f52334eb752f326aebbd3792734)) - wangjiahan

<!-- generated by git-cliff -->

---
## 0.0.1-init - 2024-07-09

Init project with pre commit hooks and GitHub Actions.

### ⛰️ Features

- build go matcher skeleton. - ([628385b](https://github.com/hedon954/go-matcher/commit/628385bfd2c9dd39813cc130215facea4a88995e)) - hedon954


