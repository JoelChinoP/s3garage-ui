import { useSearchParams } from "react-router-dom";
import { Card } from "react-daisyui";

import ObjectList from "./object-list";
import { useEffect, useState } from "react";
import ObjectListNavigator from "./object-list-navigator";
import Actions from "./actions";
import { useBucketContext } from "../context";
import ShareDialog from "./share-dialog";
import { SearchIcon, XIcon } from "lucide-react";

const getInitialPrefixes = (searchParams: URLSearchParams) => {
  const prefix = searchParams.get("prefix");
  if (prefix) {
    const paths = prefix.split("/").filter((p) => p);
    return paths.map((_, i) => paths.slice(0, i + 1).join("/") + "/");
  }
  return [];
};

const BrowseTab = () => {
  const { bucket } = useBucketContext();
  const [searchParams, setSearchParams] = useSearchParams();
  const [prefixHistory, setPrefixHistory] = useState<string[]>(
    getInitialPrefixes(searchParams)
  );
  const [curPrefix, setCurPrefix] = useState(prefixHistory.length - 1);
  const [uuidSearch, setUuidSearch] = useState("");

  useEffect(() => {
    const prefix = prefixHistory[curPrefix] || "";
    const newParams = new URLSearchParams(searchParams);
    newParams.set("prefix", prefix);
    setSearchParams(newParams);
  }, [curPrefix]);

  const gotoPrefix = (prefix: string) => {
    const history = prefixHistory.slice(0, curPrefix + 1);
    setPrefixHistory([...history, prefix]);
    setCurPrefix(history.length);
  };

  if (!bucket.keys.find((k) => k.permissions.read && k.permissions.write)) {
    return (
      <div className="p-4 min-h-[200px] flex flex-col items-center justify-center">
        <p className="text-center max-w-sm">
          You need to add a key with read & write access to your bucket to be
          able to browse it.
        </p>
      </div>
    );
  }

  return (
    <div>
      <Card className="pb-2">
        <ObjectListNavigator
          curPrefix={curPrefix}
          setCurPrefix={setCurPrefix}
          prefixHistory={prefixHistory}
          actions={<Actions prefix={prefixHistory[curPrefix] || ""} />}
        />

        <div className="px-3 pb-2">
          <div className="relative w-full max-w-xs">
            <SearchIcon
              size={14}
              className="absolute left-3 top-1/2 -translate-y-1/2 text-base-content/40 pointer-events-none"
            />
            <input
              type="text"
              placeholder="Search by UUID..."
              value={uuidSearch}
              onChange={(e) => setUuidSearch(e.target.value)}
              className="input input-sm w-full pl-8 pr-7 bg-base-200 text-sm"
            />
            {uuidSearch && (
              <button
                onClick={() => setUuidSearch("")}
                className="absolute right-2 top-1/2 -translate-y-1/2 text-base-content/40 hover:text-base-content"
              >
                <XIcon size={13} />
              </button>
            )}
          </div>
        </div>

        <ObjectList
          prefix={prefixHistory[curPrefix] || ""}
          onPrefixChange={gotoPrefix}
          uuidFilter={uuidSearch}
        />

        <ShareDialog />
      </Card>
    </div>
  );
};

export default BrowseTab;
