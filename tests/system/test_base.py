from beat-cookie-discbot import BaseTest

import os


class Test(BaseTest):

    def test_base(self):
        """
        Basic test with exiting Beat-cookie-discbot normally
        """
        self.render_config_template(
            path=os.path.abspath(self.working_dir) + "/log/*"
        )

        beat-cookie-discbot_proc = self.start_beat()
        self.wait_until(lambda: self.log_contains("beat-cookie-discbot is running"))
        exit_code = beat-cookie-discbot_proc.kill_and_wait()
        assert exit_code == 0
