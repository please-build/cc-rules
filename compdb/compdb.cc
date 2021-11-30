#include <fstream>
#include <iomanip>

#include "nlohmann/json.hpp"
#include "subprocess.hpp"

using namespace nlohmann;
typedef std::string string;

string replace(string in, const string& before, const string& after) {
  const auto idx = in.find(before);
  if (idx == string::npos) {
    return in;
  }
  return in.replace(idx, before.size(), after);
}

string trim(string in) {
  in.resize(in.find_last_not_of(" \n") + 1);
  return in;
}

int main(int argc, const char* argv[]) {
  // Get the repo root from plz (not necessarily the same as the cwd).
  auto rbuf = subprocess::check_output({"plz", "query", "reporoot"});
  const string dir = trim(string(rbuf.buf.begin(), rbuf.buf.end()));
  const string genDir = dir + "/plz-out/gen";

  auto obuf = subprocess::check_output({"plz", "query", "graph", "-c", "dbg", "--profile", "clang"});
  auto graph = json::parse(obuf.buf.begin(), obuf.buf.end());

  auto out = json::array();
  for (const auto& pkg : graph["packages"]) {
    for (const auto& target : pkg["targets"]) {
      // Checking the prefix is a pretty quick and dirty way of finding the targets
      // we consider relevant. Maybe we should check labels as well.
      if (target.contains("command") && target.contains("srcs") && target["srcs"].contains("srcs")) {
        auto cmd = target["command"].get<string>();
        if (cmd.rfind("$TOOLS_CC", 0) == 0) {  // no starts_with until C++20 :(
          // Strip the end parts where we archive the output
          auto idx = cmd.find(" && ");
          if (idx != string::npos) {
            cmd.resize(idx);
            cmd = trim(cmd);
          }
          for (const auto& src : target["srcs"]["srcs"]) {
            // Hardcode the filenames in place of variables
            string c = replace(cmd, "${SRCS_SRCS}", src);
            c = replace(c, "$TOOLS_CC", target["tools"]["cc"][0].get<string>());
            json j = {
              {"directory", genDir},
              {"command", c},
              {"file", dir + "/" + src.get<string>()},
            };
            out.emplace_back(j);
          }
        }
      }
    }
  }
  std::ofstream f;
  f.open("compile_commands.json");
  f << std::setw(4) << out << std::endl;
  f.close();
  return 0;
}
