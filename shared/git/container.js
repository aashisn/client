// @flow
import Git from '.'
import * as GitGen from '../actions/git-gen'
import * as Constants from '../constants/git'
import {anyWaiting} from '../constants/waiting'
import {compose, lifecycle, connect} from '../util/container'
import {sortBy, partition} from 'lodash-es'

const sortRepos = git => sortBy(git, ['teamname', 'name'])

const getRepos = state => {
  const git = Constants.getIdToGit(state)
  if (!git) {
    return {
      personals: [],
      teams: [],
    }
  }
  const [personals, teams] = partition(git.valueSeq().toArray(), g => !g.teamname)

  return {
    personals: sortRepos(personals).map(g => g.id),
    teams: sortRepos(teams).map(g => g.id),
  }
}

const mapStateToProps = (state, {routeState}) => {
  return {
    ...getRepos(state),
    expandedSet: routeState.get('expandedSet'),
    loading: anyWaiting(state, Constants.loadingWaitingKey),
  }
}

const mapDispatchToProps = (dispatch: any, {navigateAppend, setRouteState, routeState, navigateUp}) => ({
  _loadGit: () => dispatch(GitGen.createLoadGit()),
  onBack: () => dispatch(navigateUp()),
  onNewPersonalRepo: () => {
    dispatch(GitGen.createSetError({error: null}))
    dispatch(navigateAppend([{props: {isTeam: false}, selected: 'newRepo'}]))
  },
  onNewTeamRepo: () => {
    dispatch(GitGen.createSetError({error: null}))
    dispatch(navigateAppend([{props: {isTeam: true}, selected: 'newRepo'}]))
  },
  onShowDelete: (id: string) => {
    dispatch(GitGen.createSetError({error: null}))
    dispatch(navigateAppend([{props: {id}, selected: 'deleteRepo'}]))
  },
  onToggleExpand: (id: string) => {
    const old = routeState.get('expandedSet')
    // TODO use unique id
    setRouteState({expandedSet: old.has(id) ? old.delete(id) : old.add(id)})
  },
})

export default compose(
  connect(
    mapStateToProps,
    mapDispatchToProps,
    (s, d, o) => ({...o, ...s, ...d})
  ),
  lifecycle({
    componentDidMount() {
      this.props._loadGit()
    },
  })
)(Git)
